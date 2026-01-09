package rivo

import "context"

// TODO: tests

// TeeStream returns 2 streams that each receive a copy of each item from the input stream.
func TeeStream[T any](ctx context.Context, in Stream[T]) (Stream[T], Stream[T]) {
	ss := TeeStreamN(ctx, in, 2)
	return ss[0], ss[1]
}

// TeeStreamN returns n streams that each receive a copy of each item from the input stream.
func TeeStreamN[T any](ctx context.Context, in Stream[T], n int) []Stream[T] {
	if n <= 0 {
		panic("n must be greater than 0")
	}

	out := make([]chan T, n)
	for i := 0; i < n; i++ {
		out[i] = make(chan T)
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
