package rivo

import (
	"context"
	"sync"
)

// Parallel returns a Pipeable that applies the given pipeables to the input stream in parallel.
// The output stream will not emit any items, and it will be closed when the input stream is closed or the context is done.
func Parallel[A any](pipeables ...Pipeable[A, struct{}]) Pipeable[A, struct{}] {
	return func(ctx context.Context, in Stream[A]) Stream[struct{}] {
		out := make(chan Item[struct{}])

		go func() {
			defer close(out)

			ins := TeeN(ctx, in, len(pipeables))

			wg := sync.WaitGroup{}
			wg.Add(len(pipeables))
			defer wg.Wait()

			for i, pipeable := range pipeables {
				go func(i int, p Pipeable[A, struct{}]) {
					defer wg.Done()
					<-OrDone(ctx, p(ctx, ins[i]))
				}(i, pipeable)
			}
		}()

		return out
	}
}
