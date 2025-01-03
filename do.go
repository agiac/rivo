package rivo

import (
	"context"
	"sync"
)

// Do returns a Sync that applies the given function to each item in the stream.
// The output stream will not emit any items, and it will be closed when the input stream is closed or the context is done.
func Do[T any](f func(ctx context.Context, i Item[T]), opt ...Option) Sync[T] {
	o := mustOptions(opt...)

	return func(ctx context.Context, in Stream[T]) Stream[None] {
		out := make(chan Item[None])

		go func() {
			defer close(out)
			defer beforeClose(ctx, out, o)

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
