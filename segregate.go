package rivo

import "context"

// Segregate returns two pipelines, where the first pipeline emits items that pass the predicate, and the second pipeline emits items that do not pass the predicate.
func Segregate[T any](p Pipeline[None, T], predicate func(item Item[T]) bool) (Pipeline[None, T], Pipeline[None, T]) {
	out1 := make(chan Item[T])
	out2 := make(chan Item[T])

	p1 := func(ctx context.Context, _ Stream[None]) Stream[T] {
		return out1
	}

	p2 := func(ctx context.Context, _ Stream[None]) Stream[T] {
		return out2
	}

	go func() {
		defer close(out1)
		defer close(out2)

		// TODO: Handle context cancellation
		for item := range p(context.Background(), nil) {
			if predicate(item) {
				out1 <- item
			} else {
				out2 <- item
			}
		}
	}()

	return p1, p2
}

// SegregateErrors returns two pipelines, where the first pipeline emits items without errors, and the second pipeline emits items with errors.
func SegregateErrors[T any](p Pipeline[None, T]) (Pipeline[None, T], Pipeline[None, T]) {
	return Segregate[T](p, func(item Item[T]) bool {
		return item.Err == nil
	})
}
