package rivo

import (
	"context"
	"fmt"
)

// options contains common configuration options for a Pipeline.
type options struct {
	poolSize      int
	bufferSize    int
	stopOnError   bool
	onBeforeClose func(context.Context) error
}

var defaultOptions = options{
	poolSize:      1,
	bufferSize:    0,
	stopOnError:   false,
	onBeforeClose: nil,
}

func (o options) apply(opts ...Option) options {
	for _, opt := range opts {
		opt(&o)
	}

	return o
}

func (o options) validate() error {
	if o.poolSize < 1 {
		return fmt.Errorf("pool size must be greater than 0")
	}

	if o.bufferSize < 0 {
		return fmt.Errorf("buffer size must be greater than or equal to 0")
	}

	return nil
}

func mustOptions(opts ...Option) options {
	o := defaultOptions.apply(opts...)

	if err := o.validate(); err != nil {
		panic(fmt.Errorf("invalid options: %w", err))
	}

	return o
}

// Option is a configuration option for a Pipeline.
type Option func(*options)

// WithPoolSize sets the number of goroutines that will be used to process items. The default is 1.
func WithPoolSize(size int) Option {
	return func(o *options) {
		o.poolSize = size
	}
}

// WithBufferSize sets the size of the output channel buffer. The default is 0 (unbuffered).
func WithBufferSize(size int) Option {
	return func(o *options) {
		o.bufferSize = size
	}
}

// WithStopOnError determines whether the Pipeline should stop processing items when an error occurs. The default is false.
func WithStopOnError(stop bool) Option {
	return func(o *options) {
		o.stopOnError = stop
	}
}

// WithOnBeforeClose sets a function that will be called before the Pipeline output channel is closed.
func WithOnBeforeClose(fn func(context.Context) error) Option {
	return func(o *options) {
		// If there is already a function set, chain the new function with the existing one.
		if o.onBeforeClose != nil {
			existingFn := o.onBeforeClose
			o.onBeforeClose = func(ctx context.Context) error {
				if err := existingFn(ctx); err != nil {
					return err
				}

				return fn(ctx)
			}
		} else {
			o.onBeforeClose = fn
		}
	}
}

func beforeClose[T any](ctx context.Context, out chan<- Item[T], o options) {
	if o.onBeforeClose == nil {
		return
	}

	err := o.onBeforeClose(ctx)
	if err != nil {
		out <- Item[T]{Err: err}
	}
}
