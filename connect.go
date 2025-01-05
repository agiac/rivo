package rivo

import (
	"context"
	"sync"
)

// Connect returns a Sync that applies the given syncs to the input stream concurrently.
// The output stream will not emit any items, and it will be closed when the input stream is closed or the context is done.
func Connect[A any](pipeables ...Sync[A]) Sync[A] {
	return func(ctx context.Context, in Stream[A]) Stream[None] {
		out := make(chan Item[None])

		go func() {
			defer close(out)

			ins := TeeN(context.Background(), in, len(pipeables))

			wg := sync.WaitGroup{}
			wg.Add(len(pipeables))
			defer wg.Wait()

			for i, pipeable := range pipeables {
				go func(i int, p Sync[A]) {
					defer wg.Done()
					<-p(ctx, ins[i])
				}(i, pipeable)
			}
		}()

		return out
	}
}
