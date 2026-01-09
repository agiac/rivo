package rivo

import "context"

// Pipe pipes two workers together.
// It is a convenience function that calls Pipe2.
func Pipe[A, B, C any](a Worker[A, B], b Worker[B, C]) Worker[A, C] {
	return Pipe2(a, b)
}

// Pipe2 pipes two workers together.
func Pipe2[A, B, C any](a Worker[A, B], b Worker[B, C]) Worker[A, C] {
	return func(ctx context.Context, stream Stream[A], errs chan<- error) Stream[C] {
		return b(context.WithoutCancel(ctx), a(ctx, stream, errs), errs)
	}
}

// Pipe3 pipes three workers together.
func Pipe3[A, B, C, D any](a Worker[A, B], b Worker[B, C], c Worker[C, D]) Worker[A, D] {
	return Pipe2(Pipe2(a, b), c)
}

// Pipe4 pipes four workers together.
func Pipe4[A, B, C, D, E any](a Worker[A, B], b Worker[B, C], c Worker[C, D], d Worker[D, E]) Worker[A, E] {
	return Pipe3(Pipe2(a, b), c, d)
}

// Pipe5 pipes five workers together.
func Pipe5[A, B, C, D, E, F any](a Worker[A, B], b Worker[B, C], c Worker[C, D], d Worker[D, E], e Worker[E, F]) Worker[A, F] {
	return Pipe4(Pipe2(a, b), c, d, e)
}

// Pipe6 pipes six workers together.
func Pipe6[A, B, C, D, E, F, G any](a Worker[A, B], b Worker[B, C], c Worker[C, D], d Worker[D, E], e Worker[E, F], f Worker[F, G]) Worker[A, G] {
	return Pipe5(Pipe2(a, b), c, d, e, f)
}

// Pipe7 pipes seven workers together.
func Pipe7[A, B, C, D, E, F, G, H any](a Worker[A, B], b Worker[B, C], c Worker[C, D], d Worker[D, E], e Worker[E, F], f Worker[F, G], g Worker[G, H]) Worker[A, H] {
	return Pipe6(Pipe2(a, b), c, d, e, f, g)
}

// Pipe8 pipes eight workers together.
func Pipe8[A, B, C, D, E, F, G, H, I any](a Worker[A, B], b Worker[B, C], c Worker[C, D], d Worker[D, E], e Worker[E, F], f Worker[F, G], g Worker[G, H], h Worker[H, I]) Worker[A, I] {
	return Pipe7(Pipe2(a, b), c, d, e, f, g, h)
}

// Pipe9 pipes nine workers together.
func Pipe9[A, B, C, D, E, F, G, H, I, J any](a Worker[A, B], b Worker[B, C], c Worker[C, D], d Worker[D, E], e Worker[E, F], f Worker[F, G], g Worker[G, H], h Worker[H, I], i Worker[I, J]) Worker[A, J] {
	return Pipe8(Pipe2(a, b), c, d, e, f, g, h, i)
}

// Pipe10 pipes ten workers together.
func Pipe10[A, B, C, D, E, F, G, H, I, J, K any](a Worker[A, B], b Worker[B, C], c Worker[C, D], d Worker[D, E], e Worker[E, F], f Worker[F, G], g Worker[G, H], h Worker[H, I], i Worker[I, J], j Worker[J, K]) Worker[A, K] {
	return Pipe9(Pipe2(a, b), c, d, e, f, g, h, i, j)
}
