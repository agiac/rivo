package rivo

import "context"

// TODO: tests

// SegregateStream takes an input stream and a predicate function, and returns two streams:
// one containing items that satisfy the predicate and another containing items that do not.
func SegregateStream[T any](ctx context.Context, in Stream[T], predicate func(T) bool) (Stream[T], Stream[T]) {
	trueStream := make(chan T)
	falseStream := make(chan T)

	go func() {
		defer close(trueStream)
		defer close(falseStream)

		for item := range OrDone(ctx, in) {
			if predicate(item) {
				trueStream <- item
			} else {
				falseStream <- item
			}
		}
	}()

	return trueStream, falseStream
}

func Segregate[T any](ctx context.Context, in Stream[T], predicate func(T) bool) (Generator[T], Generator[T]) {
	trueStream, falseStream := SegregateStream(ctx, in, predicate)

	trueGen := func(ctx context.Context, _ Stream[None]) Stream[T] {
		return trueStream
	}

	falseGen := func(ctx context.Context, _ Stream[None]) Stream[T] {
		return falseStream
	}

	return trueGen, falseGen
}
