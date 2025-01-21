package rivo

import (
	"context"
	"sync"
)

// WithErrorHandler returns a pipeline that emits items from the input pipeline, and passes any errors to the error pipeline.
func WithErrorHandler[T, U any](p Pipeline[T, U], pErr Pipeline[error, None]) Pipeline[T, U] {
	return func(ctx context.Context, in Stream[T]) Stream[U] {
		out := make(chan Item[U])

		errS := make(chan Item[error])

		go func() {
			defer close(out)

			wg := sync.WaitGroup{}
			wg.Add(1)

			go func() {
				defer wg.Done()
				<-pErr(ctx, errS)
			}()

			for item := range p(ctx, in) {
				if item.Err != nil {
					errS <- Item[error]{Val: item.Err}
				} else {
					select {
					case <-ctx.Done():
						close(errS)
						return
					case out <- item:
					}
				}
			}

			close(errS)

			wg.Wait()
		}()

		return out
	}
}
