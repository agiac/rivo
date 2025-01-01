package rivo

// Item represents a single item in a data stream. It contains a value of type T and an optional error.
type Item[T any] struct {
	// Val is the value of the item when there is no error.
	Val T
	// Err is the optional error of the item.
	Err error
}

// Stream represents a data stream of items. It is a read only channel of Item[T].
type Stream[T any] <-chan Item[T]
