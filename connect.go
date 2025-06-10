package rivo

import (
	"context"
	"sync"
)

func Connect[T any](pp ...Sync[T]) Sync[T] {
	return func(ctx context.Context, in Stream[T]) Stream[None] {
		out := make(chan None)

		go func() {
			defer close(out)

			inS := TeeStreamN(ctx, in, len(pp))

			wg := sync.WaitGroup{}
			wg.Add(len(pp))

			for i, p := range pp {
				go func(i int, p Sync[T]) {
					defer wg.Done()
					<-p(ctx, inS[i])
				}(i, p)
			}

			wg.Wait()
		}()

		return out
	}
}
