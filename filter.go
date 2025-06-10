package rivo

import (
	"context"
	"fmt"
)

// Filter returns a pipeline that filters the input stream using the given function.
func Filter[T any](f func(context.Context, T) bool, opt ...FilterOption) Pipeline[T, T] {
	o := assertFilterOptions(opt)

	return ForEachOutput[T, T](
		func(ctx context.Context, val T, out chan<- T) {
			if f(ctx, val) {
				select {
				case <-ctx.Done():
					return
				case out <- val:
				}
			}
		},
		ForEachOutputPoolSize(o.poolSize),
		ForEachOutputBufferSize(o.bufferSize),
	)
}

type filterOptions struct {
	poolSize   int
	bufferSize int
}

type FilterOption func(*filterOptions) error

func FilterPoolSize(n int) FilterOption {
	return func(o *filterOptions) error {
		if n < 1 {
			return fmt.Errorf("pool size must be greater than 0")
		}

		o.poolSize = n

		return nil
	}
}

func FilterBufferSize(n int) FilterOption {
	return func(o *filterOptions) error {
		if n < 0 {
			return fmt.Errorf("buffer size must be greater than or equal to 0")
		}

		o.bufferSize = n

		return nil
	}
}

var filterDefaultOptions = filterOptions{
	poolSize:   1,
	bufferSize: 0,
}

func applyFilterOptions(opt []FilterOption) (filterOptions, error) {
	opts := filterDefaultOptions
	for _, o := range opt {
		if err := o(&opts); err != nil {
			return opts, err
		}
	}
	return opts, nil
}

func assertFilterOptions(opt []FilterOption) filterOptions {
	opts, err := applyFilterOptions(opt)
	if err != nil {
		panic(fmt.Errorf("invalid filter options: %v", err))
	}
	return opts
}
