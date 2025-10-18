package rivo

import "context"

// RunErrorSync creates and starts a new error-handling goroutine.
// It returns a channel to send errors to and a function to call to wait
// for the handler to finish.
// The returned wait function will close the error channel and wait for the
// handler goroutine to complete.
func RunErrorSync(ctx context.Context, fn func(ctx context.Context, errs <-chan error)) (chan<- error, func()) {
	done := make(chan struct{})
	errs := make(chan error, 1)

	go func() {
		defer close(done)
		fn(ctx, errs)
	}()

	return errs, func() {
		close(errs)
		<-done
	}
}

// RunErrorSyncFunc is similar to RunErrorSync but takes a function
// that handles one error at a time.
// It returns a channel to send errors to and a function to call to wait
// for the handler to finish.
// The returned wait function will close the error channel and wait for the
// handler goroutine to complete.
func RunErrorSyncFunc(ctx context.Context, fn func(ctx context.Context, err error)) (chan<- error, func()) {
	return RunErrorSync(ctx, func(ctx context.Context, errs <-chan error) {
		for err := range errs {
			fn(ctx, err)
		}
	})
}
