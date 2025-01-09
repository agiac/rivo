package rivo

import "context"

type None struct{}

// Pipeline is a function that takes a context and a stream and returns a stream. It is the building block of a data pipeline.
type Pipeline[T, U any] func(ctx context.Context, stream Stream[T]) Stream[U]

// Pipe pipes two pipeable functions together. It is a convenience function that calls Pipe2.
func Pipe[A, B, C any](a Pipeline[A, B], b Pipeline[B, C]) Pipeline[A, C] {
	return Pipe2(a, b)
}

// Pipe2 pipes two pipeable functions together.
func Pipe2[A, B, C any](a Pipeline[A, B], b Pipeline[B, C]) Pipeline[A, C] {
	return func(ctx context.Context, stream Stream[A]) Stream[C] {
		return b(context.WithoutCancel(ctx), a(ctx, stream))
	}
}

// Pipe3 pipes three pipeable functions together.
func Pipe3[A, B, C, D any](a Pipeline[A, B], b Pipeline[B, C], c Pipeline[C, D]) Pipeline[A, D] {
	return func(ctx context.Context, stream Stream[A]) Stream[D] {
		return c(context.WithoutCancel(ctx), b(context.WithoutCancel(ctx), a(ctx, stream)))
	}
}

// Pipe4 pipes four pipeable functions together.
func Pipe4[A, B, C, D, E any](a Pipeline[A, B], b Pipeline[B, C], c Pipeline[C, D], d Pipeline[D, E]) Pipeline[A, E] {
	return func(ctx context.Context, stream Stream[A]) Stream[E] {
		return d(context.WithoutCancel(ctx), c(context.WithoutCancel(ctx), b(context.WithoutCancel(ctx), a(ctx, stream))))
	}
}

// Pipe5 pipes five pipeable functions together.
func Pipe5[A, B, C, D, E, F any](a Pipeline[A, B], b Pipeline[B, C], c Pipeline[C, D], d Pipeline[D, E], e Pipeline[E, F]) Pipeline[A, F] {
	return func(ctx context.Context, stream Stream[A]) Stream[F] {
		return e(context.WithoutCancel(ctx), d(context.WithoutCancel(ctx), c(context.WithoutCancel(ctx), b(context.WithoutCancel(ctx), a(ctx, stream)))))
	}
}
