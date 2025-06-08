package rivo

import (
	"context"
	"errors"
	"github.com/agiac/rivo/core"
)

var ErrEOS = errors.New("end of stream")

// FromFunc returns a Generator that emits items returned by the given function.
// The returned stream will emit items until the function returns ErrEOS.
// Error items are emitted if the function returns an error other than ErrEOS.
func FromFunc[T any](f func(context.Context) (T, error), options ...core.FromFuncOption) Generator[T] {
	return core.FromFunc[Item[T]](func(ctx context.Context) (Item[T], bool) {
		v, err := f(ctx)
		if errors.Is(err, ErrEOS) {
			return Item[T]{}, false // Signal end of stream
		}

		if err != nil {
			return Item[T]{Err: err}, true // Emit error item
		}

		return Item[T]{Val: v}, true // Emit value item
	}, options...)
}
