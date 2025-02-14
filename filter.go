package rivo

import (
	"context"
	"sync"
)

// Filter returns a pipeline that filters the input stream using the given function.
func Filter[T any](f func(context.Context, Item[T]) (bool, error), opt ...Option) Pipeline[T, T] {
	o := mustOptions(opt...)

	return func(ctx context.Context, stream Stream[T]) Stream[T] {
		out := make(chan Item[T], o.bufferSize)

		wg := sync.WaitGroup{}
		wg.Add(o.poolSize)

		go func() {
			defer close(out)
			defer beforeClose(ctx, out, o)

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
								if o.stopOnError {
									return
								}
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
