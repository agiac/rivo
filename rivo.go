// Package rivo is a library for stream processing.
package rivo

import "context"

// None is a type that represents no value.
// It is typically used as the input type of generator pipeline that does not
// depend on any input stream or for a sync pipeline that does not emit any items.
type None struct{}

// Item represents a single item in a data stream. It contains a value of type T and an optional error.
type Item[T any] struct {
	// Val is the value of the item when there is no error.
	Val T
	// Err is the optional error of the item.
	Err error
}

// Stream represents a data stream of items. It is a read only channel of Item[T].
type Stream[T any] <-chan Item[T]

// Pipeline is a function that takes a context and a stream and returns a stream of the same type or a different type.
type Pipeline[T, U any] func(ctx context.Context, stream Stream[T]) Stream[U]
