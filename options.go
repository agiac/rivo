package rivo

import (
	"fmt"
)

// options contains common configuration options for a Pipeable.
type options struct {
	poolSize    int
	bufferSize  int
	stopOnError bool
}

var defaultOptions = options{
	poolSize:    1,
	bufferSize:  0,
	stopOnError: false,
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

// Option is a configuration option for a Pipeable.
type Option func(*options)

// WithPoolSize sets the number of goroutines that will be used to process items concurrently. The default is 1.
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

// WithStopOnError determines whether the Pipeable should stop processing items when an error occurs. The default is false.
func WithStopOnError(stop bool) Option {
	return func(o *options) {
		o.stopOnError = stop
	}
}
