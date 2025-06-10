package rivo

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// TODO: tests

// ForEachOutput returns a pipeline that applies a function to each item from the input stream.
// The function can write directly to the output channel. The output channel should not be closed by the function,
// since the output stream will be closed when the input stream is closed or the context is done.
// ForEachOutput panics if invalid options are provided.
func ForEachOutput[T, U any](f func(ctx context.Context, val T, out chan<- U), opt ...ForEachOutputOption) Pipeline[T, U] {
	o := mustForEachOutputOptions(opt)

	return func(ctx context.Context, in Stream[T]) Stream[U] {
		out := make(chan U, o.bufferSize)

		go func() {
			defer close(out)
			defer o.onBeforeClose(ctx)

			wg := sync.WaitGroup{}
			wg.Add(o.poolSize)

			for i := 0; i < o.poolSize; i++ {
				go func() {
					defer wg.Done()

					for {
						select {
						case <-ctx.Done():
							return
						case v, ok := <-in:
							if !ok {
								return
							}

							f(ctx, v, out)
						}
					}
				}()
			}

			wg.Wait()
		}()

		return out
	}
}

type forEachOutputOptions struct {
	poolSize      int
	bufferSize    int
	onBeforeClose func(context.Context)
}

type ForEachOutputOption func(*forEachOutputOptions) error

func ForEachOutputPoolSize(poolSize int) ForEachOutputOption {
	return func(o *forEachOutputOptions) error {
		if poolSize < 1 {
			return errors.New("poolSize must be greater than 0")
		}
		o.poolSize = poolSize
		return nil
	}
}

func ForEachOutputBufferSize(bufferSize int) ForEachOutputOption {
	return func(o *forEachOutputOptions) error {
		if bufferSize < 0 {
			return errors.New("bufferSize must be greater than or equal to 0")
		}
		o.bufferSize = bufferSize
		return nil
	}
}

func ForEachOutputOnBeforeClose(f func(context.Context)) ForEachOutputOption {
	return func(o *forEachOutputOptions) error {
		if f == nil {
			return errors.New("onBeforeClose must not be nil")
		}
		o.onBeforeClose = f
		return nil
	}
}

func newDefaultForEachOutputOptions() *forEachOutputOptions {
	return &forEachOutputOptions{
		poolSize:      1,
		bufferSize:    0,
		onBeforeClose: func(ctx context.Context) {},
	}
}

func applyForEachOutputOptions(opts []ForEachOutputOption) (*forEachOutputOptions, error) {
	o := newDefaultForEachOutputOptions()
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	return o, nil
}

func mustForEachOutputOptions(opts []ForEachOutputOption) *forEachOutputOptions {
	o, err := applyForEachOutputOptions(opts)
	if err != nil {
		panic(fmt.Sprintf("invalid ForEachOutputOption: %v", err))
	}
	return o
}
