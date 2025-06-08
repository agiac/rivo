package rivo

import (
	"github.com/agiac/rivo/core"
)

// Connect returns a sync pipelines that applies the given syncs pipelines to the input stream concurrently.
// The output stream will not emit any items, and it will be closed when the input stream is closed or the context is done.
func Connect[A any](pp ...Sync[A]) Sync[A] {
	return core.Connect[Item[A]](pp...)
}
