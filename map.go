package rivo

import (
	"context"
	"fmt"
)

// Map returns a pipeline that applies a function to each item from the input stream.
func Map[T, U any](f func(context.Context, T) (U, error), opt ...MapOption) Pipeline[T, U] {
	o := mustMapOptions(opt)

	return ForEachOutput[T, U](
		func(ctx context.Context, val T, out chan<- U, errs chan<- error) {
			v, err := f(ctx, val)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errs <- err:
				}
				return
			}

			select {
			case <-ctx.Done():
				return
			case out <- v:
			}
		},
		ForEachOutputPoolSize(o.poolSize),
		ForEachOutputBufferSize(o.bufferSize),
	)
}

type mapOptions struct {
	poolSize   int
	bufferSize int
}

type MapOption func(*mapOptions) error

func MapPoolSize(poolSize int) MapOption {
	return func(o *mapOptions) error {
		if poolSize < 1 {
			return fmt.Errorf("poolSize must be greater than 0")
		}
		o.poolSize = poolSize
		return nil
	}
}

func MapBufferSize(bufferSize int) MapOption {
	return func(o *mapOptions) error {
		if bufferSize < 0 {
			return fmt.Errorf("bufferSize must be greater than or equal to 0")
		}
		o.bufferSize = bufferSize
		return nil
	}
}

func newDefaultMapOptions() *mapOptions {
	return &mapOptions{
		poolSize:   1,
		bufferSize: 0,
	}
}

func applyMapOptions(opts []MapOption) (*mapOptions, error) {
	o := newDefaultMapOptions()

	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	return o, nil
}

func mustMapOptions(opts []MapOption) *mapOptions {
	o, err := applyMapOptions(opts)
	if err != nil {
		panic(fmt.Sprintf("invalid MapOption: %v", err))
	}
	return o
}
