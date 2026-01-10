package rivo

import "context"

// Fanout connects a Worker to multiple Sync workers.
// It takes a Worker and a variable number of Sync workers,
// and returns a new Sync worker that first processes input items
// using the provided Worker, and then fans out the output to all
// the provided Sync workers for further processing.
func Fanout[T, U any](g Worker[T, U], ss ...Sync[U]) Sync[T] {
	return func(ctx context.Context, in <-chan T, errs chan<- error) <-chan None {
		return Connect[U](ss...)(ctx, g(ctx, in, errs), errs)
	}
}
