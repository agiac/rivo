package core

import (
	"context"
)

// Flatten returns a Pipeline that flattens a Stream of slices into a Stream of individual items.
func Flatten[T any]() Pipeline[[]T, T] {
	return func(ctx context.Context, in Stream[[]T]) Stream[T] {
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

					for _, val := range item {
						select {
						case <-ctx.Done():
							return
						case out <- val:
						}
					}
				}
			}
		}()

		return out
	}
}
