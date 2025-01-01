package rivo

import "context"

// Pipeable is a function that takes a context and a stream and returns a stream. It is the building block of a data pipeline.
type Pipeable[T, U any] func(ctx context.Context, stream Stream[T]) Stream[U]

// Pipe pipes two pipeable functions together. It is a convenience function that calls Pipe2.
func Pipe[A, B, C any](a Pipeable[A, B], b Pipeable[B, C]) Pipeable[A, C] {
	return Pipe2(a, b)
}

// Pipe2 pipes two pipeable functions together.
func Pipe2[A, B, C any](a Pipeable[A, B], b Pipeable[B, C]) Pipeable[A, C] {
	return func(ctx context.Context, stream Stream[A]) Stream[C] {
		return b(ctx, a(ctx, stream))
	}
}

// Pipe3 pipes three pipeable functions together.
func Pipe3[A, B, C, D any](a Pipeable[A, B], b Pipeable[B, C], c Pipeable[C, D]) Pipeable[A, D] {
	return func(ctx context.Context, stream Stream[A]) Stream[D] {
		return c(ctx, b(ctx, a(ctx, stream)))
	}
}

// Pipe4 pipes four pipeable functions together.
func Pipe4[A, B, C, D, E any](a Pipeable[A, B], b Pipeable[B, C], c Pipeable[C, D], d Pipeable[D, E]) Pipeable[A, E] {
	return func(ctx context.Context, stream Stream[A]) Stream[E] {
		return d(ctx, c(ctx, b(ctx, a(ctx, stream))))
	}
}

// Pipe5 pipes five pipeable functions together.
func Pipe5[A, B, C, D, E, F any](a Pipeable[A, B], b Pipeable[B, C], c Pipeable[C, D], d Pipeable[D, E], e Pipeable[E, F]) Pipeable[A, F] {
	return func(ctx context.Context, stream Stream[A]) Stream[F] {
		return e(ctx, d(ctx, c(ctx, b(ctx, a(ctx, stream)))))
	}
}
