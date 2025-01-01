package rivo

import "context"

// Pipeable is a function that takes a context and a stream and returns a stream. It is the building block of a data pipeline.
type Pipeable[T, U any] func(ctx context.Context, stream Stream[T]) Stream[U]

// Pipeline represents a data pipeline that start with a generator Pipeable and can be extended with more Pipeable functions.
type Pipeline[T, U any] Pipeable[T, U]

// Collect collects all items from a Pipeline and returns them as a slice.
func (p Pipeline[T, U]) Collect() []Item[U] {
	return Collect(p(context.Background(), nil))
}

// CollectWithContext collects all items from a Pipeline and returns them as a slice. If the context is cancelled, it stops collecting items.
func (p Pipeline[T, U]) CollectWithContext(ctx context.Context) []Item[U] {
	return CollectWithContext(ctx, p(ctx, nil))
}

// Pipe pipes two Pipeable together and returns a Pipeline. It is a convenience function that calls Pipe2.
func Pipe[A, B, C any](a Pipeable[A, B], b Pipeable[B, C]) Pipeline[A, C] {
	return Pipe2(a, b)
}

// Pipe2 pipes two Pipeable together and returns a Pipeline.
func Pipe2[A, B, C any](a Pipeable[A, B], b Pipeable[B, C]) Pipeline[A, C] {
	return func(ctx context.Context, stream Stream[A]) Stream[C] {
		return b(ctx, a(ctx, stream))
	}
}

// Pipe3 pipes three Pipeable together and returns a Pipeline.
func Pipe3[A, B, C, D any](a Pipeable[A, B], b Pipeable[B, C], c Pipeable[C, D]) Pipeline[A, D] {
	return func(ctx context.Context, stream Stream[A]) Stream[D] {
		return c(ctx, b(ctx, a(ctx, stream)))
	}
}

// Pipe4 pipes four Pipeable together and returns a Pipeline.
func Pipe4[A, B, C, D, E any](a Pipeable[A, B], b Pipeable[B, C], c Pipeable[C, D], d Pipeable[D, E]) Pipeline[A, E] {
	return func(ctx context.Context, stream Stream[A]) Stream[E] {
		return d(ctx, c(ctx, b(ctx, a(ctx, stream))))
	}
}

// Pipe5 pipes five Pipeable together and returns a Pipeline.
func Pipe5[A, B, C, D, E, F any](a Pipeable[A, B], b Pipeable[B, C], c Pipeable[C, D], d Pipeable[D, E], e Pipeable[E, F]) Pipeline[A, F] {
	return func(ctx context.Context, stream Stream[A]) Stream[F] {
		return e(ctx, d(ctx, c(ctx, b(ctx, a(ctx, stream)))))
	}
}
