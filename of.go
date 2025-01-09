package rivo

import "context"

// Of returns a generator Pipeline that emits the given items. The input stream is ignored.
func Of[T any](items ...T) Pipeline[None, T] {
	return func(ctx context.Context, _ Stream[None]) Stream[T] {
		out := make(chan Item[T])

		go func() {
			defer close(out)

			for _, item := range items {
				select {
				case <-ctx.Done():
					out <- Item[T]{Err: ctx.Err()}
					return
				case out <- Item[T]{Val: item}:
				}
			}
		}()

		return out
	}
}
