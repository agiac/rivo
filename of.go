package rivo

import (
	"context"
	"github.com/agiac/rivo/core"
)

// Of returns a generator Pipeline that emits the given items. The input stream is ignored.
func Of[T any](items ...T) Generator[T] {
	return func(ctx context.Context, stream core.Stream[core.None]) core.Stream[Item[T]] {
		out := make(chan Item[T])

		go func() {
			defer close(out)

			for _, item := range items {
				select {
				case <-ctx.Done():
					return
				case out <- Item[T]{Val: item}:
				}
			}
		}()

		return out
	}
}
