package rivo

import "context"

// Fork returns two pipeables that apply the given pipeable to the input stream.
func Fork[T, U any](p Pipeable[T, U]) (Pipeable[T, U], Pipeable[T, U]) {
	ps := ForkN(p, 2)
	return ps[0], ps[1]
}

// ForkN returns n pipeables that apply the given pipeable to the input stream.
func ForkN[T, U any](p Pipeable[T, U], n int) []Pipeable[T, U] {
	ps := make([]Pipeable[T, U], n)

	for i := 0; i < n; i++ {
		ps[i] = func(ctx context.Context, in Stream[T]) Stream[U] {
			return p(ctx, in)
		}
	}

	return ps
}
