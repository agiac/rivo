package rivo

import "context"

type ForEachFunc[T any] = func(context.Context, Item[T]) error

// ForEach returns a pipeline that applies a function to each item from the input stream.
// It is intended for side effect and the output stream will only emit the errors returned by the function.
func ForEach[T any](f ForEachFunc[T], opt ...Option) Pipeable[T, struct{}] {
	forEach := Map[T, struct{}](func(ctx context.Context, item Item[T]) (struct{}, error) {
		return struct{}{}, f(ctx, item)
	}, opt...)

	filterNoErrors := Filter(func(ctx context.Context, i Item[struct{}]) (bool, error) {
		return i.Err != nil, nil
	})

	return Pipe(forEach, filterNoErrors)
}
