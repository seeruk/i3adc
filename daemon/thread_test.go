package daemon

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBackgroundThread(t *testing.T) {
	t.Run("should start and stop successfully under normal circumstances", func(t *testing.T) {
		ctx, cfn := context.WithCancel(context.Background())
		defer cfn()

		thread := newTestThread(nil, nil)
		errCh := NewBackgroundThread(ctx, thread)

		select {
		case <-errCh:
			t.Error("expected background thread to still be running")
		default:
		}

		cfn()

		timeout := time.NewTimer(5 * time.Second)

		select {
		case <-errCh:
		case <-timeout.C:
			t.Error("expected background thread to no longer be running")
		}

		assert.Equal(t, 1, thread.startCalls)
		assert.Equal(t, 1, thread.stopCalls)
	})

	t.Run("should not call stop if start failed", func(t *testing.T) {
		ctx, cfn := context.WithCancel(context.Background())
		defer cfn()

		thread := newTestThread(errors.New("oops"), nil)
		errCh := NewBackgroundThread(ctx, thread)

		timeout := time.NewTimer(5 * time.Second)

		select {
		case <-errCh:
		case <-timeout.C:
			t.Error("expected background thread to return error")
		}

		assert.Equal(t, 1, thread.startCalls)
		assert.Equal(t, 0, thread.stopCalls)
	})

	t.Run("should return the error from stopping, if one is returned", func(t *testing.T) {
		ctx, cfn := context.WithCancel(context.Background())
		defer cfn()

		thread := newTestThread(nil, errors.New("oops"))
		errCh := NewBackgroundThread(ctx, thread)

		select {
		case <-errCh:
			t.Error("expected background thread to still be running")
		default:
		}

		cfn()

		timeout := time.NewTimer(5 * time.Second)

		select {
		case <-errCh:
		case <-timeout.C:
			t.Error("expected background thread to return error")
		}

		assert.Equal(t, 1, thread.startCalls)
		assert.Equal(t, 1, thread.stopCalls)
	})
}

// testThread is a Thread used to test the functionality around threads.
type testThread struct {
	done chan struct{}

	startError error
	stopError  error

	startCalls int
	stopCalls  int
}

// newTestThread returns a new test thread instance.
func newTestThread(startError, stopError error) *testThread {
	return &testThread{
		done:       make(chan struct{}),
		startError: startError,
		stopError:  stopError,
	}
}

// Start is a "mock" thread start function.
func (t *testThread) Start() error {
	t.startCalls++

	if t.startError == nil {
		<-t.done
	}

	return t.startError
}

// Stop is a "mock" thread stop function.
func (t *testThread) Stop() error {
	t.stopCalls++
	t.done <- struct{}{}

	return t.stopError
}
