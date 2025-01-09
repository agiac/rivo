package rivo

import "context"

// Segregate returns a function that returns two pipeables, where the first pipeable emits items that pass the predicate, and the second pipeable emits items that do not pass the predicate.
func Segregate[T, U any](p Pipeline[T, U], predicate func(ctx context.Context, item Item[U]) bool) func(context.Context, Stream[T]) (Pipeline[None, U], Pipeline[None, U]) {
	return func(ctx context.Context, in Stream[T]) (Pipeline[None, U], Pipeline[None, U]) {
		out1 := make(chan Item[U])
		out2 := make(chan Item[U])

		p1 := func(ctx context.Context, _ Stream[None]) Stream[U] {
			return out1
		}

		p2 := func(ctx context.Context, _ Stream[None]) Stream[U] {
			return out2
		}

		go func() {
			defer close(out1)
			defer close(out2)

			for item := range p(ctx, in) {
				if predicate(ctx, item) {
					select {
					case out1 <- item:
					case <-ctx.Done():
						return
					}
				} else {
					select {
					case out2 <- item:
					case <-ctx.Done():
						return
					}
				}
			}
		}()

		return p1, p2
	}
}
