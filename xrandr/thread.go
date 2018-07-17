package xrandr

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/seeruk/i3adc/event"
	"github.com/seeruk/i3adc/logging"
	"github.com/seeruk/i3adc/state"
)

// Thread is a process that will wait for events from an event channel, and based on those events,
// will attempt to trigger updates to the display configuration.
type Thread struct {
	ctx     context.Context
	cfn     context.CancelFunc
	backend state.Backend
	client  *Client
	logger  logging.Logger
	eventCh <-chan event.Event
}

// NewThread returns a new output thread instance.
func NewThread(backend state.Backend, client *Client, logger logging.Logger, eventCh <-chan event.Event) *Thread {
	logger = logger.With("module", "xrandr/thread")

	return &Thread{
		backend: backend,
		client:  client,
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
		case evt := <-t.eventCh:
			err := t.onEvent(evt)
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

func (t *Thread) onEvent(evt event.Event) error {
	t.logger.Debug("event occurred")

	currentLayout, err := t.client.GetOutputs()
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

		var lastXPos int

		// For all connected outputs, activate the preferred mode, and do a little bit of auto
		// config to do with things like rotation, reflection, and position. Fairly standard.
		for _, output := range currentLayout {
			var args []string

			if !output.IsConnected {
				// We have to turn off ones that shouldn't be there!
				args = []string{
					"--output", output.Name,
					"--off",
				}
			} else {
				var preferredMode *Mode
				for _, mode := range output.Modes {
					m := mode
					if mode.IsPreferred {
						preferredMode = &m
					}
				}

				// Otherwise, set some sensible default settings.
				args = []string{
					"--output", output.Name,
					"--pos", fmt.Sprintf("%dx0", lastXPos),
					"--rotate", "normal",
					"--reflect", "normal",
				}

				if lastXPos == 0 {
					args = append(args, "--primary")
				}

				if preferredMode == nil {
					args = append(args, "--auto")
				} else {
					args = append(args, "--mode", preferredMode.Name)
					lastXPos += int(preferredMode.Width)
				}
			}

			t.logger.Debugw("running command", "command", "xrandr", "args", args)

			cmd := exec.Command("xrandr", args...)

			// Attempt to run xrandr command...
			err = cmd.Run()
			if err != nil {
				return err
			}
		}

		// Re-fetch layout, so our changes are applied to our in-memory representation.
		currentLayout, err := t.client.GetOutputs()
		if err != nil {
			return err
		}

		currentLayoutBS, err := json.Marshal(currentLayout)
		if err != nil {
			return err
		}

		t.backend.Write(state.KeyLatestLayout, []byte(hash))
		t.backend.Write(hash, currentLayoutBS)
	case !evt.IsStartup && hash == latestHash:
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

		var savedLayout []Output

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
					"--mode", fmt.Sprintf("%dx%d", output.Width, output.Height),
					"--pos", fmt.Sprintf("%dx%d", output.OffsetX, output.OffsetY),
					"--rotate", output.Rotation.String(),
					"--reflect", output.Reflection.String(),
				}

				if output.IsPrimary {
					args = append(args, "--primary")
				}
			}

			t.logger.Debugw("running command", "command", "xrandr", "args", args)

			cmd := exec.Command("xrandr", args...)

			// Attempt to run xrandr command...
			err = cmd.Run()
			if err != nil {
				return err
			}
		}

		t.backend.Write(state.KeyLatestLayout, []byte(hash))
	}

	return nil
}
