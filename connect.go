package rivo

import (
	"context"
	"sync"
)

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

func Fanout[T, U any](g Worker[T, U], ss ...Sync[U]) Sync[T] {
	return func(ctx context.Context, in <-chan T, errs chan<- error) <-chan None {
		return Connect[U](ss...)(ctx, g(ctx, in, errs), errs)
	}
}
