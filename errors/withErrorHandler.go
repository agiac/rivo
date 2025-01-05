package errors

import (
	"context"
	"sync"

	"github.com/agiac/rivo"
)

// WithErrorHandler returns a pipeable that connects the input pipeable to an error handling pipeable.
// The items that don't contain errors are passed to the output stream, while the items that contain errors are passed to the error handling pipeable.
func WithErrorHandler[T, U any](p rivo.Pipeable[T, U], errHandler rivo.Sync[struct{}]) rivo.Pipeable[T, U] {
	return func(ctx context.Context, in rivo.Stream[T]) rivo.Stream[U] {
		out := make(chan rivo.Item[U])

		vals, errs := rivo.Segregate(p, func(ctx context.Context, item rivo.Item[U]) bool {
			return item.Err == nil
		})(ctx, in)

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			<-rivo.Pipe3(
				errs,
				rivo.Map(func(ctx context.Context, i rivo.Item[U]) (struct{}, error) {
					return struct{}{}, i.Err
				}),
				errHandler,
			)(ctx, nil)
		}()

		go func() {
			defer close(out)
			for i := range vals(ctx, nil) {
				select {
				case <-ctx.Done():
				case out <- i:
				}
			}

			wg.Wait()
		}()

		return out
	}
}
