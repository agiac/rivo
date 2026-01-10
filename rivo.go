package rivo

import (
	"context"
)

// None is a type that represents no value.
// It is typically used as the input type of generator worker that does not
// depend on any input channel or for a sync worker that does not emit any items.
type None struct{}

// Worker is the core abstraction representing a processing unit that
// takes input of type T from a channel, processes it, and emits output of type U to another channel.
// It also receives a context for cancellation and a channel for reporting errors.
type Worker[T, U any] = func(ctx context.Context, in <-chan T, errs chan<- error) <-chan U

// Generator is a worker that generates items of type T without any input.
type Generator[T any] = Worker[None, T]

// Sync is a worker that processes items of type T and does not emit any items.
type Sync[T any] = Worker[T, None]
