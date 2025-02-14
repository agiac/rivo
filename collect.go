package rivo

import (
	"context"
)

// Collect collects all items from the stream and returns them as a slice.
func Collect[T any](in Stream[T]) []Item[T] {
	return CollectWithContext(context.Background(), in)
}

// CollectWithContext collects all items from the stream and returns them as a slice. If the context is cancelled, it stops collecting items.
func CollectWithContext[T any](ctx context.Context, in Stream[T]) []Item[T] {
	var items []Item[T]

	for item := range OrDone(ctx, in) {
		items = append(items, item)
	}

	return items
}
