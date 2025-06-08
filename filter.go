package rivo

import (
	"context"
	"github.com/agiac/rivo/core"
)

// Filter returns a pipeline that filters the input stream using the given function.
// The function should return true for items that should be included in the output stream.
// Errors in the input stream are always propagated, and the filter function is not called for those items.
func Filter[T any](f func(context.Context, T) bool, opt ...core.FilterOption) Pipeline[T, T] {
	return core.Filter[Item[T]](func(ctx context.Context, i Item[T]) bool {
		if i.Err != nil {
			return true // Always propagate errors
		}

		return f(ctx, i.Val)
	}, opt...)
}
