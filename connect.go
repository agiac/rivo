package rivo

import (
	"context"
	"sync"
)

// Connect returns a sync pipelines that applies the given syncs pipelines to the input stream concurrently.
// The output stream will not emit any items, and it will be closed when the input stream is closed or the context is done.
func Connect[A any](pp ...Pipeline[A, None]) Pipeline[A, None] {
	return func(ctx context.Context, in Stream[A]) Stream[None] {
		out := make(chan Item[None])

		go func() {
			defer close(out)

			ins := TeeN(context.Background(), in, len(pp))

			wg := sync.WaitGroup{}
			wg.Add(len(pp))
			defer wg.Wait()

			for i, p := range pp {
				go func(i int, p Pipeline[A, None]) {
					defer wg.Done()
					<-p(ctx, ins[i])
				}(i, p)
			}
		}()

		return out
	}
}
