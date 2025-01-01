package rivo

import (
	"fmt"
)

type options struct {
	poolSize        int
	bufferSize      int
	stopOnError     bool
	extraValidation func(o options) error
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

	if o.extraValidation != nil {
		return o.extraValidation(o)
	}

	return nil
}

func mustOptions(opts ...Option) options {
	o := defaultOptions().apply(opts...)

	if err := o.validate(); err != nil {
		panic(fmt.Errorf("invalid options: %w", err))
	}

	return o
}

func defaultOptions() options {
	return options{
		poolSize:    1,
		bufferSize:  0,
		stopOnError: false,
	}
}

type Option func(*options)

func withExtraValidation(f func(o options) error) Option {
	return func(o *options) {
		o.extraValidation = f
	}
}

func WithPoolSize(size int) Option {
	return func(o *options) {
		o.poolSize = size
	}
}

func WithBufferSize(size int) Option {
	return func(o *options) {
		o.bufferSize = size
	}
}

func WithStopOnError(stop bool) Option {
	return func(o *options) {
		o.stopOnError = stop
	}
}
