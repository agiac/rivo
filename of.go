package rivo

import "context"

// Of returns a Generator that emits the given items.
func Of[T any](items ...T) Generator[T] {
	return func(ctx context.Context, _ Stream[None], _ chan<- error) Stream[T] {
		out := make(chan T)

		go func() {
			defer close(out)

			for _, item := range items {
				select {
				case <-ctx.Done():
					return
				case out <- item:
				}
			}
		}()

		return out
	}
}
