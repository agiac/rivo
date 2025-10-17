package rivo

import (
	"context"
)

// Flatten returns a Pipeline that flattens a Stream of slices into a Stream of individual items.
func Flatten[T any]() Pipeline[[]T, T] {
	return ForEachOutput[[]T, T](func(ctx context.Context, val []T, out chan<- T, errs chan<- error) {
		for _, item := range val {
			select {
			case <-ctx.Done():
				return
			case out <- item:
			}
		}
	})
}
