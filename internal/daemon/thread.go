package daemon

import "context"

// Thread is a generic interface for some sort of process that can be started and stopped. It is not
// necessarily run in the background, and it may not necessarily have to be stopped by calling stop
// (i.e. it could end on it's own). A thread should honour it's stop method and bail gracefully, but
// cleanly.
type Thread interface {
	// Start starts some work. The work may end when Stop is called, or it may end on it's own.
	Start() error
	// Stop attempts to stop a started thread.
	Stop() error
}

// NewBackgroundThread starts the given Thread, and then waits for the given context to signal that
// it should stop it's work, allowing the thread to (hopefully) end gracefully. The returned channel
// will receive a message when the thread has finished stopping. It should only receive one message.
func NewBackgroundThread(ctx context.Context, thread Thread) <-chan error {
	bail := make(chan error, 1) // bail stops pointless work and leaking goroutines.
	done := make(chan error, 1) // done is the channel the result is sent through.

	// This goroutine watches the context, waiting for a cancellation signal. If the thread ends
	// itself before a cancellation is received from the context, then it will be signaled to
	// "bail", preventing the goroutine from leaking.
	go func() {
		select {
		case <-ctx.Done():
			// If we get signalled to stop, try stop.
			bail <- thread.Stop()
		case <-bail:
			// If the other routine signals us to stop, we've already stopped.
		}
	}()

	// This goroutine actually starts the thread. The thread should then block until it is either
	// cancelled, or it's work is done. If the thread stops on it's own, this goroutine will signal
	// the one above to "bail". A result is then sent down the done channel to notify where this
	// thread is being used that it has ended.
	go func() {
		err := thread.Start()

		select {
		case err = <-bail:
			// We'll hit this if stop was called and returned an error very quickly.
		default:
			// We'll hit this if stop hasn't been called, in which case, we stop the other routine.
			bail <- nil
		}

		done <- err
	}()

	return done
}
