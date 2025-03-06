package rivo

import (
	"context"
)

// Tee returns two pipelines that each receive a copy of each item from the input stream.
func Tee[None, T any](p Pipeline[None, T]) (Pipeline[None, T], Pipeline[None, T]) {
	tees := TeeN(p, 2)
	return tees[0], tees[1]
}

// TeeN returns n pipelines that each receive a copy of each item from the input stream.
func TeeN[None, T any](p Pipeline[None, T], n int) []Pipeline[None, T] {
	if n <= 1 {
		panic("n must be greater than 1")
	}

	streams := teeStream[T](p(context.Background(), nil), n)

	pipes := make([]Pipeline[None, T], n)
	for i := 0; i < n; i++ {
		tee := streams[i]
		pipes[i] = func(ctx context.Context, _ Stream[None]) Stream[T] {
			return tee
		}
	}

	return pipes
}

// teeStream returns n streams that each receive a copy of each item from the input stream.
func teeStream[T any](in Stream[T], n int) []Stream[T] {
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

		for item := range in {
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
