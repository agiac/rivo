package rivo

import (
	"context"
	"errors"
	"sync"
)

func FromFunc[T any](f func(ctx context.Context) (T, error), options ...Option) Pipeable[struct{}, T] {
	o := mustOptions(options...)

	return func(ctx context.Context, stream Stream[struct{}]) Stream[T] {
		out := make(chan Item[T], o.bufferSize)

		//onFinish := func() {
		//	if err := o.onFinish(ctx); err != nil {
		//		out <- Item[T]{Err: err}
		//	}
		//}

		go func() {
			defer close(out)
			//defer onFinish()

			wg := sync.WaitGroup{}
			wg.Add(o.poolSize)

			for i := 0; i < o.poolSize; i++ {
				go func() {
					defer wg.Done()

					for {
						v, err := f(ctx)
						if errors.Is(err, ErrEOS) {
							return
						}

						select {
						case <-ctx.Done():
							out <- Item[T]{Err: ctx.Err()}
							return
						default:
							out <- Item[T]{Val: v, Err: err}
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
