package xrandr

import (
	"context"

	"github.com/seeruk/i3adc/logging"
	"github.com/seeruk/i3adc/state"
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

	props, err := getProps()
	if err != nil {
		return err
	}

	outputs, err := parseProps(props)
	if err != nil {
		return err
	}

	hash, err := calculateHashForOutputs(outputs)
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
		// for all connected outputs. The user can then set their configuration themselves to update
		// the saved configuration. This is a new layout.
		t.logger.Infow("creating a new configuration", "hash", hash)
	case hash == latestHash:
		// If the hash is the same, we want to update the existing layout at that hash. Either this
		// output configuration has been used before, or the user has just updated it. Technically,
		// all we need to do is that update here...
		t.logger.Infow("updating an existing configuration", "hash", hash)
	default:
		// Otherwise, we aren't updating layout, or creating a new one, we're simply switching to
		// another layout. In other words, we should just apply the `savedLayoutBS` configuration.
		t.logger.Infow("switching to existing configuration", "hash", hash, "previous_hash", latestHash)
	}

	return nil
}
