package rivo

import (
	"context"
	"sync"
)

// Connect connects multiple Sync workers in parallel.
// It takes a variable number of Sync workers and returns a new Sync worker
// that runs all the provided workers concurrently on the same input channel.
// The output channel of the returned Sync worker will be closed once all
// the connected workers have completed their processing.
func Connect[T any](pp ...Sync[T]) Sync[T] {
	return func(ctx context.Context, in <-chan T, errs chan<- error) <-chan None {
		out := make(chan None)

		go func() {
			defer close(out)

			inS := TeeN(ctx, in, len(pp))

			wg := sync.WaitGroup{}
			wg.Add(len(pp))

			for i, p := range pp {
				go func(i int, p Sync[T]) {
					defer wg.Done()
					<-p(ctx, inS[i], errs)
				}(i, p)
			}

			wg.Wait()
		}()

		return out
	}
}
