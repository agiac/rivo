package rivo

import (
	"context"
)

// TeeN returns n streams that each receive a copy of each item from the input stream.
func TeeN[T any](ctx context.Context, in Stream[T], n int) []Stream[T] {
	if n <= 1 {
		panic("n must be greater than 1")
	}

	out := make([]chan Item[T], n)
	for i := 0; i < n; i++ {
		out[i] = make(chan Item[T])
	}

	go func() {
		defer func() {
			for i := 0; i < n; i++ {
				close(out[i])
			}
		}()

		for item := range OrDone(ctx, in) {
			for i := 0; i < n; i++ {
				out[i] <- item
			}
		}
	}()

	streams := make([]Stream[T], n)
	for i := 0; i < n; i++ {
		streams[i] = out[i]
	}

	return streams
}

func Tee2[T, U any](p Pipeline[T, U]) func(context.Context, Stream[T]) (Pipeline[None, U], Pipeline[None, U]) {
	return func(ctx context.Context, s Stream[T]) (Pipeline[None, U], Pipeline[None, U]) {
		res := Tee2N[T, U](p, 2)(ctx, s)
		return res[0], res[1]
	}
}

func Tee2N[T, U any](p Pipeline[T, U], n int) func(context.Context, Stream[T]) []Pipeline[None, U] {
	if n <= 1 {
		panic("n must be greater than 1")
	}

	return func(ctx context.Context, s Stream[T]) []Pipeline[None, U] {
		out := make([]chan Item[U], n)
		for i := 0; i < n; i++ {
			out[i] = make(chan Item[U])
		}

		go func() {
			defer func() {
				for i := 0; i < n; i++ {
					close(out[i])
				}
			}()

			for item := range OrDone(ctx, p(ctx, s)) {
				for i := 0; i < n; i++ {
					select {
					case out[i] <- item:
					case <-ctx.Done():
						out[i] <- Item[U]{Err: ctx.Err()}
					}
				}
			}
		}()

		pipes := make([]Pipeline[None, U], n)
		for i := 0; i < n; i++ {
			out := out[i]
			pipes[i] = func(ctx context.Context, _ Stream[None]) Stream[U] {
				return out
			}
		}

		return pipes
	}
}
