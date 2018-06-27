package xrandr

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/seeruk/i3adc/logging"
	"github.com/seeruk/i3adc/state"
	"github.com/seeruk/i3adc/xrandr/props"
)

// Thread is a process that will wait for events from an event channel, and based on those events,
// will attempt to trigger updates to the display configuration.
type Thread struct {
	ctx     context.Context
	cfn     context.CancelFunc
	backend state.Backend
	logger  logging.Logger
	eventCh <-chan struct{}
}

// NewThread returns a new output thread instance.
func NewThread(backend state.Backend, logger logging.Logger, eventCh <-chan struct{}) *Thread {
	logger = logger.With("module", "xrandr/thread")

	return &Thread{
		backend: backend,
		eventCh: eventCh,
		logger:  logger,
	}
}

// Start begins waiting for events in the event channel. When an event occurs, this thread will
// trigger behaviour to update the display configuration, if necessary.
func (t *Thread) Start() error {
	t.logger.Info("thread started")
	t.ctx, t.cfn = context.WithCancel(context.Background())

	for {
		select {
		case <-t.ctx.Done():
			t.logger.Info("thread stopped")
			return t.ctx.Err()
		case <-t.eventCh:
			err := t.onEvent()
			if err != nil {
				t.logger.Errorw("error handling event",
					"error", err.Error(),
				)
			}
		}
	}
}

// Stop attempts to stop this thread.
func (t *Thread) Stop() error {
	t.logger.Infow("thread stopping")

	if t.ctx != nil && t.cfn != nil {
		t.cfn()
	}

	return nil
}

func (t *Thread) onEvent() error {
	t.logger.Debug("event occurred")

	currentLayout, err := getLayout()
	if err != nil {
		return err
	}

	hash, err := calculateHashForOutputs(currentLayout)
	if err != nil {
		return err
	}

	t.logger.Debugw("calculated hash", "hash", hash)

	latestHashBS, err := t.backend.Read(state.KeyLatestLayout)
	if err != nil {
		return err
	}

	latestHash := string(latestHashBS)

	t.logger.Debugw("latest hash", "hash", latestHash)

	savedLayoutBS, err := t.backend.Read(hash)
	if err != nil {
		return err
	}

	switch {
	case savedLayoutBS == nil:
		// If we haven't got a layout stored for this hash, we should activate the preferred mode
		// for all connected currentLayout. The user can then set their configuration themselves to update
		// the saved configuration. This is a new layout.
		t.logger.Infow("creating a new configuration", "hash", hash)

		// TODO(seeruk): Move logic away from here...

		var lastXPos int

		// For all connected outputs, activate the preferred mode, and do a little bit of auto
		// config to do with things like rotation, reflection, and position. Fairly standard.
		for _, output := range currentLayout {
			var args []string

			if !output.IsConnected || !output.IsEnabled {
				// We have to turn off ones that shouldn't be there!
				args = []string{
					"--output", output.Name,
					"--off",
				}
			} else {
				// Otherwise, set some sensible default settings.
				args = []string{
					"--output", output.Name,
					"--auto",
					"--pos", fmt.Sprintf("%dx0", lastXPos),
					"--rotate", "normal",
					"--reflect", "normal",
				}

				if lastXPos == 0 {
					args = append(args, "--primary")
				}

				// NOTE(seeruk): Not here, but later when we're working on generating xrandr commands,
				// we will potentially need to take into account rotation. If a screen is on it's side,
				// is it's height it's width in terms of positioning?
				lastXPos = lastXPos + int(output.Resolution.Width)
			}

			cmd := exec.Command("xrandr", args...)

			// Attempt to run xrandr command...
			err = cmd.Run()
			if err != nil {
				return err
			}
		}

		currentLayout, err := getLayout()
		if err != nil {
			return err
		}

		currentLayoutBS, err := json.Marshal(currentLayout)
		if err != nil {
			return err
		}

		t.backend.Write(state.KeyLatestLayout, []byte(hash))
		t.backend.Write(hash, currentLayoutBS)
	case hash == latestHash:
		// If the hash is the same, we want to update the existing layout at that hash. Either this
		// output configuration has been used before, or the user has just updated it. Technically,
		// all we need to do is that update here...
		t.logger.Infow("updating an existing configuration", "hash", hash)

		currentLayoutBS, err := json.Marshal(currentLayout)
		if err != nil {
			return err
		}

		t.backend.Write(state.KeyLatestLayout, []byte(hash))
		t.backend.Write(hash, currentLayoutBS)
	default:
		// Otherwise, we aren't updating layout, or creating a new one, we're simply switching to
		// another layout. In other words, we should just apply the `savedLayoutBS` configuration.
		t.logger.Infow("switching to existing configuration", "hash", hash, "previous_hash", latestHash)

		var savedLayout []props.Output

		err := json.Unmarshal(savedLayoutBS, &savedLayout)
		if err != nil {
			return err
		}

		for _, output := range savedLayout {
			var args []string

			if !output.IsConnected || !output.IsEnabled {
				args = []string{
					"--output", output.Name,
					"--off",
				}
			} else {
				args = []string{
					"--output", output.Name,
					"--mode", fmt.Sprintf("%dx%d", output.Resolution.Width, output.Resolution.Height),
					"--pos", fmt.Sprintf("%dx%d", output.Position.OffsetX, output.Position.OffsetY),
					"--rotate", "normal", // @TODO
					"--reflect", "normal", // @TODO
				}

				if output.IsPrimary {
					args = append(args, "--primary")
				}
			}

			cmd := exec.Command("xrandr", args...)

			// Attempt to run xrandr command...
			err = cmd.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getLayout() ([]props.Output, error) {
	ps, err := getProps()
	if err != nil {
		return []props.Output{}, err
	}

	return parseProps(ps)
}
