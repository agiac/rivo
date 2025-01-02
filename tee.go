package rivo

import (
	"context"
)

// Tee returns two streams that each receive a copy of each item from the input stream. It is equivalent to TeeN(ctx, in, 2).
func Tee[T any](ctx context.Context, in Stream[T]) (Stream[T], Stream[T]) {
	out := TeeN(ctx, in, 2)
	return out[0], out[1]
}

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
