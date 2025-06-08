package core

import (
	"context"
)

// None is a type that represents no value.
// It is typically used as the input type of generator pipeline that does not
// depend on any input stream or for a sync pipeline that does not emit any items.
type None struct{}

// Stream represents a data stream of items. It is a read only channel of type T.
type Stream[T any] <-chan T

// Pipeline is a function that takes a context and a stream and returns a stream of the same type or a different type.
type Pipeline[T, U any] func(ctx context.Context, stream Stream[T]) Stream[U]

// Generator is a pipeline that generates items of type T without any input.
type Generator[T any] = Pipeline[None, T]

// Sync is a pipeline that processes items of type T and does not emit any items.
type Sync[T any] = Pipeline[T, None]
