package core

import (
	"context"
	"fmt"
	"sync"
)

// Map returns a pipeline that applies a function to each item from the input stream.
func Map[T, U any](f func(context.Context, T) U, opt ...MapOption) Pipeline[T, U] {
	o := mustMapOptions(opt)

	return func(ctx context.Context, stream Stream[T]) Stream[U] {
		out := make(chan U, o.bufferSize)

		wg := sync.WaitGroup{}
		wg.Add(o.poolSize)

		go func() {
			defer close(out)

			for range o.poolSize {
				go func() {
					defer wg.Done()

					for item := range OrDone(ctx, stream) {
						v := f(ctx, item)

						select {
						case <-ctx.Done():
							return
						case out <- v:
						}
					}
				}()
			}

			wg.Wait()
		}()

		return out
	}
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
