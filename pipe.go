package rivo

import "context"

// Pipe pipes two pipelines together.
// It is a convenience function that calls Pipe2.
func Pipe[A, B, C any](a Pipeline[A, B], b Pipeline[B, C]) Pipeline[A, C] {
	return Pipe2(a, b)
}

// Pipe2 pipes two pipelines together.
func Pipe2[A, B, C any](a Pipeline[A, B], b Pipeline[B, C]) Pipeline[A, C] {
	return func(ctx context.Context, stream Stream[A]) Stream[C] {
		return b(context.WithoutCancel(ctx), a(ctx, stream))
	}
}

// Pipe3 pipes three pipelines together.
func Pipe3[A, B, C, D any](a Pipeline[A, B], b Pipeline[B, C], c Pipeline[C, D]) Pipeline[A, D] {
	return Pipe2(Pipe2(a, b), c)
}

// Pipe4 pipes four pipelines together.
func Pipe4[A, B, C, D, E any](a Pipeline[A, B], b Pipeline[B, C], c Pipeline[C, D], d Pipeline[D, E]) Pipeline[A, E] {
	return Pipe3(Pipe2(a, b), c, d)
}

// Pipe5 pipes five pipelines together.
func Pipe5[A, B, C, D, E, F any](a Pipeline[A, B], b Pipeline[B, C], c Pipeline[C, D], d Pipeline[D, E], e Pipeline[E, F]) Pipeline[A, F] {
	return Pipe4(Pipe2(a, b), c, d, e)
}
