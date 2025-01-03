package rivo

import (
	"context"
	"errors"
	"sync"
)

var ErrEOS = errors.New("end of stream")

// FromFunc returns a generator Pipeable that emits items generated by the given function. The input stream is ignored.
// The returned stream will emit items until the function returns ErrEOS.
func FromFunc[T any](f func(ctx context.Context) (T, error), options ...Option) Generator[T] {
	o := mustOptions(options...)

	return func(ctx context.Context, stream Stream[None]) Stream[T] {
		out := make(chan Item[T], o.bufferSize)

		go func() {
			defer close(out)
			defer beforeClose(ctx, out, o)

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
