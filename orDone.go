package rivo

import "context"

func OrDone[T any](ctx context.Context, in Stream[T]) Stream[T] {
	out := make(chan Item[T])

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				out <- Item[T]{Err: ctx.Err()}
				return
			case item, ok := <-in:
				if !ok {
					return
				}

				select {
				case out <- item:
				case <-ctx.Done():
					out <- Item[T]{Err: ctx.Err()}
					return
				}
			}
		}
	}()

	return out
}
