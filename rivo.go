package rivo

import (
	"context"
)

// None is a type that represents no value.
// It is typically used as the input type of generator worker that does not
// depend on any input stream or for a sync worker that does not emit any items.
type None struct{}

// Stream represents a data stream of items. It is a read only channel of type T.
type Stream[T any] = <-chan T

// Worker is a function that takes a context and a stream and returns a stream of the same type or a different type.
type Worker[T, U any] = func(ctx context.Context, stream Stream[T], errs chan<- error) Stream[U]

// Generator is a worker that generates items of type T without any input.
type Generator[T any] = Worker[None, T]

// Sync is a worker that processes items of type T and does not emit any items.
type Sync[T any] = Worker[T, None]
