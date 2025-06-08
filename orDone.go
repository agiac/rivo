package rivo

import (
	"context"
	"github.com/agiac/rivo/core"
)

// OrDone is a utility function that returns a channel that will be closed when the context is done.
func OrDone[T any](ctx context.Context, in Stream[T]) Stream[T] {
	return core.OrDone[Item[T]](ctx, in)
}
