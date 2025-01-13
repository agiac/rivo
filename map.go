package rivo

import (
	"context"
	"sync"
)

type MapFunc[T, U any] = func(context.Context, Item[T]) (U, error)

// Map returns a pipeline that applies a function to each item from the input stream.
func Map[T, U any](f MapFunc[T, U], opt ...Option) Pipeline[T, U] {
	o := mustOptions(opt...)

	return func(ctx context.Context, stream Stream[T]) Stream[U] {
		out := make(chan Item[U], o.bufferSize)

		wg := sync.WaitGroup{}
		wg.Add(o.poolSize)

		go func() {
			defer close(out)
			defer beforeClose(ctx, out, o)

			for range o.poolSize {
				go func() {
					defer wg.Done()

					for item := range OrDone(ctx, stream) {
						v, err := f(ctx, item)

						select {
						case <-ctx.Done():
							out <- Item[U]{Err: ctx.Err()}
							return
						case out <- Item[U]{Val: v, Err: err}:
							if err != nil && o.stopOnError {
								return
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
