package rivo

import "context"

// Fanin combines multiple Generator workers into a single Generator.
// It takes a variable number of Generator workers and returns a new Generator
// that merges the outputs of all the provided generators into a single output channel.
func Fanin[T any](gg ...Generator[T]) Generator[T] {
	return func(ctx context.Context, _ <-chan None, errs chan<- error) <-chan T {
		ins := make([]<-chan T, len(gg))

		for i, g := range gg {
			ins[i] = g(ctx, nil, errs)
		}

		return Merge(ctx, ins...)
	}
}
