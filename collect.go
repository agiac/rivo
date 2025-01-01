package rivo

import (
	"context"
)

func Collect[T any](in Stream[T]) []Item[T] {
	return CollectWithContext(context.Background(), in)
}

func CollectWithContext[T any](ctx context.Context, in Stream[T]) []Item[T] {
	var items []Item[T]

	for item := range OrDone(ctx, in) {
		items = append(items, item)
	}

	return items
}
