package core

import (
	"context"
	"fmt"
)

// Do returns a sync pipeline that applies the given function to each item in the stream.
// The output stream will not emit any items, and it will be closed when the input stream is closed or the context is done.
func Do[T any](f func(context.Context, T), opt ...DoOption) Sync[T] {
	o := assertDoOptions(opt)

	return Sync[T](ForEachOutput[T, None](
		func(ctx context.Context, val T, out chan<- None) {
			f(ctx, val)
		},
		ForEachOutputPoolSize(o.poolSize),
		ForEachOutputOnBeforeClose(o.onBeforeClose),
	))

	//return func(ctx context.Context, in Stream[T]) Stream[None] {
	//	out := make(chan None)
	//
	//	go func() {
	//		defer close(out)
	//
	//		wg := sync.WaitGroup{}
	//		wg.Add(o.poolSize)
	//
	//		for i := 0; i < o.poolSize; i++ {
	//			go func() {
	//				defer wg.Done()
	//
	//				for item := range OrDone(ctx, in) {
	//					f(ctx, item)
	//				}
	//			}()
	//		}
	//
	//		wg.Wait()
	//	}()
	//
	//	return out
	//}
}

type doOptions struct {
	poolSize      int
	onBeforeClose func(context.Context)
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

func DoOnBeforeClose(fn func(context.Context)) DoOption {
	return func(o *doOptions) error {
		if fn == nil {
			return fmt.Errorf("onBeforeClose function cannot be nil")
		}

		o.onBeforeClose = fn

		return nil
	}
}

func newDefaultDoOptions() *doOptions {
	return &doOptions{
		poolSize:      1,
		onBeforeClose: func(ctx context.Context) {},
	}
}

func applyDoOptions(opt []DoOption) (*doOptions, error) {
	opts := newDefaultDoOptions()
	for _, o := range opt {
		if err := o(opts); err != nil {
			return opts, err
		}
	}
	return opts, nil
}

func assertDoOptions(opt []DoOption) *doOptions {
	opts, err := applyDoOptions(opt)
	if err != nil {
		panic(fmt.Sprintf("invalid Do options: %v", err))
	}
	return opts
}
