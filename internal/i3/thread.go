package i3

import (
	"context"

	"github.com/seeruk/i3adc/internal/logging"
	"go.i3wm.org/i3"
)

// Thread is a background thread designed to push output events from the i3 IPC into a channel to
// trigger other functionality in i3adc.
type Thread struct {
	ctx    context.Context
	cfn    context.CancelFunc
	logger logging.Logger
	msgCh  chan<- struct{}
	rcvr   *i3.EventReceiver
}

// NewThread creates a new output event thread instance.
func NewThread(logger logging.Logger, msgCh chan<- struct{}) *Thread {
	logger = logger.With("module", "internal/i3/thread")

	return &Thread{
		logger: logger,
		msgCh:  msgCh,
	}
}

// Start begins waiting for events from i3, pushing them onto the message channel when possible.
func (t *Thread) Start() error {
	t.logger.Info("thread started")
	t.ctx, t.cfn = context.WithCancel(context.Background())

	t.rcvr = i3.Subscribe(i3.OutputEventType)
	defer t.rcvr.Close()

	// Use a goroutine to allow this thread to be stopped. This goroutine will not die though, which
	// is very unfortunate, but shouldn't be a problem for i3adc, given it's current implementation.
	go func() {
		for t.rcvr.Next() {
			// Check context here so that we break out of the loop if possible. It may not always be
			// the case, meaning sometimes we may have a routine being leaked, at least for 5
			// seconds. This should only really happen when we're quitting anyway.
			select {
			case <-t.ctx.Done():
				break
			default:
			}

			t.msgCh <- struct{}{}
		}
	}()

	// Wait for stop signal.
	<-t.ctx.Done()

	t.logger.Info("thread stopped")

	return nil
}

// Stop attempts to stop this thread.
func (t *Thread) Stop() error {
	t.logger.Infow("thread stopping")

	if t.ctx != nil && t.cfn != nil {
		t.cfn()
	}

	return nil
}
