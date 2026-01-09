package rivo

import "context"

// TODO: tests

// SegregateChan takes an input channel and a predicate function, and returns two channels:
// one containing items that satisfy the predicate and another containing items that do not.
func SegregateChan[T any](ctx context.Context, in <-chan T, predicate func(T) bool) (<-chan T, <-chan T) {
	trueCh := make(chan T)
	falseCh := make(chan T)

	go func() {
		defer close(trueCh)
		defer close(falseCh)

		for item := range OrDone(ctx, in) {
			if predicate(item) {
				trueCh <- item
			} else {
				falseCh <- item
			}
		}
	}()

	return trueCh, falseCh
}

func Segregate[T any](ctx context.Context, in <-chan T, predicate func(T) bool) (Generator[T], Generator[T]) {
	trueCh, falseCh := SegregateChan(ctx, in, predicate)

	trueGen := func(ctx context.Context, _ <-chan None, errs chan<- error) <-chan T {
		return trueCh
	}

	falseGen := func(ctx context.Context, _ <-chan None, errs chan<- error) <-chan T {
		return falseCh
	}

	return trueGen, falseGen
}
