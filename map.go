package rivo

import (
	"context"
	"github.com/agiac/rivo/core"
)

// Map returns a pipeline that applies a function to each item from the input stream.
func Map[T, U any](f func(context.Context, Item[T]) (U, error), opt ...core.MapOption) Pipeline[T, U] {
	return core.Map[Item[T], Item[U]](func(ctx context.Context, i Item[T]) Item[U] {
		if i.Err != nil {
			return Item[U]{Err: i.Err}
		}

		val, err := f(ctx, i)
		if err != nil {
			return Item[U]{Err: err}
		}

		return Item[U]{Val: val}
	}, opt...)
}
