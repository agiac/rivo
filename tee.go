package rivo

import (
	"context"
)

// Tee returns a pipeline that applies the given pipeline to the input stream and returns two pipelines that each receive a copy of each item from the input stream.
func Tee[T, U any](p Pipeline[T, U]) func(context.Context, Stream[T]) (Pipeline[None, U], Pipeline[None, U]) {
	return func(ctx context.Context, s Stream[T]) (Pipeline[None, U], Pipeline[None, U]) {
		res := TeeN[T, U](p, 2)(ctx, s)
		return res[0], res[1]
	}
}

// TeeN returns a pipeline that applies the given pipeline to the input stream and returns n pipelines that each receive a copy of each item from the input stream.
func TeeN[T, U any](p Pipeline[T, U], n int) func(context.Context, Stream[T]) []Pipeline[None, U] {
	if n <= 1 {
		panic("n must be greater than 1")
	}

	return func(ctx context.Context, s Stream[T]) []Pipeline[None, U] {
		teeS := teeStream(ctx, s, n)

		pipes := make([]Pipeline[None, U], n)
		for i := 0; i < n; i++ {
			tee := teeS[i]
			pipes[i] = func(ctx context.Context, _ Stream[None]) Stream[U] {
				return p(ctx, tee)
			}
		}

		return pipes
	}
}

// teeStream returns n streams that each receive a copy of each item from the input stream.
func teeStream[T any](ctx context.Context, in Stream[T], n int) []Stream[T] {
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
