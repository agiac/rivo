package core

import "context"

// TODO: tests

// TeeStream returns 2 streams that each receive a copy of each item from the input stream.
func TeeStream[T any](ctx context.Context, in Stream[T]) (Stream[T], Stream[T]) {
	ss := TeeStreamN(ctx, in, 2)
	return ss[0], ss[1]
}

// TeeStreamN returns n streams that each receive a copy of each item from the input stream.
func TeeStreamN[T any](ctx context.Context, in Stream[T], n int) []Stream[T] {
	if n <= 1 {
		panic("n must be greater than 1")
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

// Tee returns 2 generators that each receive a copy of each item from the input stream.
func Tee[T any](ctx context.Context, in Stream[T]) (Generator[T], Generator[T]) {
	streams := TeeStreamN(ctx, in, 2)

	gen1 := func(ctx context.Context, _ Stream[None]) Stream[T] {
		return streams[0]
	}

	gen2 := func(ctx context.Context, _ Stream[None]) Stream[T] {
		return streams[1]
	}

	return gen1, gen2
}

// TeeN returns n generators that each receive a copy of each item from the input stream.
func TeeN[T any](ctx context.Context, in Stream[T], n int) []Generator[T] {
	if n <= 1 {
		panic("n must be greater than 1")
	}

	streams := TeeStreamN(ctx, in, n)
	generators := make([]Generator[T], n)

	for i := 0; i < n; i++ {
		generators[i] = func(ctx context.Context, _ Stream[None]) Stream[T] {
			return streams[i]
		}
	}

	return generators
}
