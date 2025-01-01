package rivo

import "context"

func Of[T any](items ...T) Pipeable[struct{}, T] {
	return func(ctx context.Context, stream Stream[struct{}]) Stream[T] {
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
