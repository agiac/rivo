package rivo

import (
	"context"
	"fmt"
	"sync"
)

type doOptions struct {
	poolSize int
}

type DoOption func(*doOptions) error

func DoPoolSize(n int) DoOption {
	return func(o *doOptions) error {
		if n < 1 {
			return fmt.Errorf("pool size must be greater than 0")
		}

		o.poolSize = n

		return nil
	}
}

var doDefaultOptions = doOptions{
	poolSize: 1,
}

func applyDoOptions(opt []DoOption) (doOptions, error) {
	opts := doDefaultOptions
	for _, o := range opt {
		if err := o(&opts); err != nil {
			return opts, err
		}
	}
	return opts, nil
}

// Do returns a sync pipeline that applies the given function to each item in the stream.
// The output stream will not emit any items, and it will be closed when the input stream is closed or the context is done.
func Do[T any](f func(context.Context, Item[T]), opt ...DoOption) Pipeline[T, None] {
	o, err := applyDoOptions(opt)
	if err != nil {
		panic(fmt.Errorf("invalid Do options: %v", err))
	}

	return func(ctx context.Context, in Stream[T]) Stream[None] {
		out := make(chan Item[None])

		go func() {
			defer close(out)

			wg := sync.WaitGroup{}
			wg.Add(o.poolSize)

			for i := 0; i < o.poolSize; i++ {
				go func() {
					defer wg.Done()

					for item := range OrDone(ctx, in) {
						f(ctx, item)
					}
				}()
			}

			wg.Wait()
		}()

		return out
	}
}
