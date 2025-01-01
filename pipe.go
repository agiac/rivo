package rivo

import "context"

type Pipeable[T, U any] func(ctx context.Context, stream Stream[T]) Stream[U]

func (p Pipeable[T, U]) Collect() []Item[U] {
	return Collect(p(context.Background(), nil))
}

func (p Pipeable[T, U]) CollectWithContext(ctx context.Context) []Item[U] {
	return CollectWithContext(ctx, p(ctx, nil))
}

func Pipe[A, B, C any](a Pipeable[A, B], b Pipeable[B, C]) Pipeable[A, C] {
	return Pipe2(a, b)
}

func Pipe2[A, B, C any](a Pipeable[A, B], b Pipeable[B, C]) Pipeable[A, C] {
	return func(ctx context.Context, stream Stream[A]) Stream[C] {
		return b(ctx, a(ctx, stream))
	}
}

func Pipe3[A, B, C, D any](a Pipeable[A, B], b Pipeable[B, C], c Pipeable[C, D]) Pipeable[A, D] {
	return func(ctx context.Context, stream Stream[A]) Stream[D] {
		return c(ctx, b(ctx, a(ctx, stream)))
	}
}

func Pipe4[A, B, C, D, E any](a Pipeable[A, B], b Pipeable[B, C], c Pipeable[C, D], d Pipeable[D, E]) Pipeable[A, E] {
	return func(ctx context.Context, stream Stream[A]) Stream[E] {
		return d(ctx, c(ctx, b(ctx, a(ctx, stream))))
	}
}

func Pipe5[A, B, C, D, E, F any](a Pipeable[A, B], b Pipeable[B, C], c Pipeable[C, D], d Pipeable[D, E], e Pipeable[E, F]) Pipeable[A, F] {
	return func(ctx context.Context, stream Stream[A]) Stream[F] {
		return e(ctx, d(ctx, c(ctx, b(ctx, a(ctx, stream)))))
	}
}
