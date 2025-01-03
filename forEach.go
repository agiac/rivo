package rivo

import "context"

// ForEach returns a Transformer that applies a function to each item from the input stream.
// It is intended for side effect and the output stream will only emit the errors returned by the function.
func ForEach[T any](f func(context.Context, Item[T]) error, opt ...Option) Transformer[T, struct{}] {
	forEach := Map[T, struct{}](func(ctx context.Context, item Item[T]) (struct{}, error) {
		return struct{}{}, f(ctx, item)
	}, opt...)

	filterNoErrors := Filter(func(ctx context.Context, i Item[struct{}]) (bool, error) {
		return i.Err != nil, nil
	})

	return Pipe(forEach, filterNoErrors)
}
