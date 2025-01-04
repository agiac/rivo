package rivo

import "context"

type None struct{}

// Pipeable is a function that takes a context and a stream and returns a stream. It is the building block of a data pipeline.
type Pipeable[T, U any] func(ctx context.Context, stream Stream[T]) Stream[U]

// Generator is a pipeable function that does not read from its input stream. It starts a new stream from scratch.
type Generator[T any] = Pipeable[None, T]

// Sync is a pipeable function that does not emit any items. It is used at the end of a pipeline.
type Sync[T any] = Pipeable[T, None]

// Transformer is a pipeable that reads from its input stream and emits items to its output stream.
type Transformer[T, U any] = Pipeable[T, U]

// Pipe pipes two pipeable functions together. It is a convenience function that calls Pipe2.
func Pipe[A, B, C any](a Pipeable[A, B], b Pipeable[B, C]) Pipeable[A, C] {
	return Pipe2(a, b)
}

// Pipe2 pipes two pipeable functions together.
func Pipe2[A, B, C any](a Pipeable[A, B], b Pipeable[B, C]) Pipeable[A, C] {
	return func(ctx context.Context, stream Stream[A]) Stream[C] {
		return b(context.WithoutCancel(ctx), a(ctx, stream))
	}
}

// Pipe3 pipes three pipeable functions together.
func Pipe3[A, B, C, D any](a Pipeable[A, B], b Pipeable[B, C], c Pipeable[C, D]) Pipeable[A, D] {
	return func(ctx context.Context, stream Stream[A]) Stream[D] {
		return c(context.WithoutCancel(ctx), b(context.WithoutCancel(ctx), a(ctx, stream)))
	}
}

// Pipe4 pipes four pipeable functions together.
func Pipe4[A, B, C, D, E any](a Pipeable[A, B], b Pipeable[B, C], c Pipeable[C, D], d Pipeable[D, E]) Pipeable[A, E] {
	return func(ctx context.Context, stream Stream[A]) Stream[E] {
		return d(context.WithoutCancel(ctx), c(context.WithoutCancel(ctx), b(context.WithoutCancel(ctx), a(ctx, stream))))
	}
}

// Pipe5 pipes five pipeable functions together.
func Pipe5[A, B, C, D, E, F any](a Pipeable[A, B], b Pipeable[B, C], c Pipeable[C, D], d Pipeable[D, E], e Pipeable[E, F]) Pipeable[A, F] {
	return func(ctx context.Context, stream Stream[A]) Stream[F] {
		return e(context.WithoutCancel(ctx), d(context.WithoutCancel(ctx), c(context.WithoutCancel(ctx), b(context.WithoutCancel(ctx), a(ctx, stream)))))
	}
}
