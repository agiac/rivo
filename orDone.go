package rivo

import "context"

// OrDone is a utility function that returns a channel that will be closed when the context is done.
func OrDone[T any](ctx context.Context, in Stream[T]) Stream[T] {
	out := make(chan T)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-in:
				if !ok {
					return
				}

				select {
				case out <- item:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out
}
