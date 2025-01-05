package errors

import (
	"context"
	"sync"

	"github.com/agiac/rivo"
)

func WithErrorHandler[T, U any](sVal rivo.Pipeable[T, U], sErr rivo.Sync[T]) rivo.Pipeable[T, U] {
	fErr := rivo.Filter(func(ctx context.Context, i rivo.Item[T]) (bool, error) {
		return i.Err != nil, nil
	})

	fVal := rivo.Filter(func(ctx context.Context, i rivo.Item[T]) (bool, error) {
		return i.Err == nil, nil
	})

	pVal := rivo.Pipe(fVal, sVal)
	pErr := rivo.Pipe(fErr, sErr)

	return func(ctx context.Context, in rivo.Stream[T]) rivo.Stream[U] {
		out := make(chan rivo.Item[U])

		go func() {
			defer close(out)

			inVal, inErr := rivo.Tee(context.Background(), in)

			wg := sync.WaitGroup{}
			wg.Add(2)
			defer wg.Wait()

			go func() {
				defer wg.Done()
				for i := range pVal(ctx, inVal) {
					select {
					case <-ctx.Done():
					case out <- i:
					}
				}
			}()

			go func() {
				defer wg.Done()
				select {
				case <-ctx.Done():
				case <-pErr(ctx, inErr):
				}
			}()
		}()

		return out
	}
}

func WithErrorHandler2[T, U any](p rivo.Pipeable[T, U], errHandler rivo.Sync[U]) rivo.Pipeable[T, U] {
	return func(ctx context.Context, in rivo.Stream[T]) rivo.Stream[U] {
		out := make(chan rivo.Item[U])

		vals, errs := rivo.Segregate(p, func(ctx context.Context, item rivo.Item[U]) bool {
			return item.Err == nil
		})(ctx, in)

		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			<-rivo.Pipe(errs, errHandler)(ctx, nil)
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
