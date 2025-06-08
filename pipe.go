package rivo

import (
	"context"
)

func Pipe[A, B, C any](a Pipeline[A, B], b Pipeline[B, C]) Pipeline[A, C] {
	return Pipe2(a, b)
}

func Pipe2[A, B, C any](a Pipeline[A, B], b Pipeline[B, C]) Pipeline[A, C] {
	return func(ctx context.Context, stream Stream[A]) Stream[C] {
		return b(context.WithoutCancel(ctx), a(ctx, stream))
	}
}

func Pipe3[A, B, C, D any](a Pipeline[A, B], b Pipeline[B, C], c Pipeline[C, D]) Pipeline[A, D] {
	return Pipe2(Pipe2(a, b), c)
}

func Pipe4[A, B, C, D, E any](a Pipeline[A, B], b Pipeline[B, C], c Pipeline[C, D], d Pipeline[D, E]) Pipeline[A, E] {
	return Pipe3(Pipe2(a, b), c, d)
}

func Pipe5[A, B, C, D, E, F any](a Pipeline[A, B], b Pipeline[B, C], c Pipeline[C, D], d Pipeline[D, E], e Pipeline[E, F]) Pipeline[A, F] {
	return Pipe4(Pipe2(a, b), c, d, e)
}
