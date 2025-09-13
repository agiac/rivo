package rivo

import (
	"context"
	"fmt"
)

// FilterMap returns a pipeline that filters and maps items from the input stream.
func FilterMap[T, U any](f func(context.Context, T) (bool, U), opt ...FilterMapOption) Pipeline[T, U] {
	o := assertFilterMapOptions(opt)

	return ForEachOutput[T, U](
		func(ctx context.Context, val T, out chan<- U) {
			keep, mapped := f(ctx, val)
			if !keep {
				return
			}

			select {
			case <-ctx.Done():
				return
			case out <- mapped:
			}
		},
		ForEachOutputPoolSize(o.poolSize),
		ForEachOutputBufferSize(o.bufferSize),
	)
}

type filterMapOptions struct {
	poolSize   int
	bufferSize int
}

type FilterMapOption func(*filterMapOptions) error

func FilterMapPoolSize(n int) FilterMapOption {
	return func(o *filterMapOptions) error {
		if n < 1 {
			return fmt.Errorf("pool size must be greater than 0")
		}

		o.poolSize = n

		return nil
	}
}

func FilterMapBufferSize(n int) FilterMapOption {
	return func(o *filterMapOptions) error {
		if n < 0 {
			return fmt.Errorf("buffer size must be greater than or equal to 0")
		}

		o.bufferSize = n

		return nil
	}
}

var filterMapDefaultOptions = filterMapOptions{
	poolSize:   1,
	bufferSize: 0,
}

func applyFilterMapOptions(opt []FilterMapOption) (filterMapOptions, error) {
	opts := filterMapDefaultOptions
	for _, o := range opt {
		if err := o(&opts); err != nil {
			return opts, err
		}
	}
	return opts, nil
}

func assertFilterMapOptions(opt []FilterMapOption) filterMapOptions {
	opts, err := applyFilterMapOptions(opt)
	if err != nil {
		panic(fmt.Errorf("invalid filter options: %v", err))
	}
	return opts
}
