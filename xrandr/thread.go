package xrandr

import (
	"context"

	"github.com/seeruk/i3adc/logging"
)

// Thread is a process that will wait for events from an event channel, and based on those events,
// will attempt to trigger updates to the display configuration.
type Thread struct {
	ctx     context.Context
	cfn     context.CancelFunc
	logger  logging.Logger
	eventCh <-chan struct{}
}

// NewThread returns a new output thread instance.
func NewThread(logger logging.Logger, eventCh <-chan struct{}) *Thread {
	logger = logger.With("module", "xrandr/thread")

	return &Thread{
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

	return nil
}