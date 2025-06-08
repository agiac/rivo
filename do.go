package rivo

import (
	"context"
	"github.com/agiac/rivo/core"
)

// Do returns a sync pipeline that applies the given function to each item in the stream.
// The output stream will not emit any items, and it will be closed when the input stream is closed or the context is done.
func Do[T any](f func(context.Context, Item[T]), opt ...core.DoOption) Sync[T] {
	return core.Do[Item[T]](f, opt...)
}
