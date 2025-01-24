package rivo

import (
	"context"
	"fmt"
	"sync"
)

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

// Filter returns a pipeline that filters the input stream using the given function.
func Filter[T any](f func(context.Context, Item[T]) (bool, error), opt ...FilterOption) Pipeline[T, T] {
	o, err := applyFilterOptions(opt)
	if err != nil {
		panic(fmt.Errorf("invalid Filter options: %v", err))
	}

	return func(ctx context.Context, stream Stream[T]) Stream[T] {
		out := make(chan Item[T], o.bufferSize)

		wg := sync.WaitGroup{}
		wg.Add(o.poolSize)

		go func() {
			defer close(out)

			for range o.poolSize {
				go func() {
					defer wg.Done()

					for item := range OrDone(ctx, stream) {
						ok, err := f(ctx, item)

						select {
						case <-ctx.Done():
							out <- Item[T]{Err: ctx.Err()}
							return
						default:
							if err != nil {
								out <- Item[T]{Err: err}
							}

							if ok {
								out <- item
							}
						}
					}
				}()
			}

			wg.Wait()
		}()

		return out
	}
}
