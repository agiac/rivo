package rivo

import (
	"context"
	"time"
)

func Batch[T any](n int, maxWait time.Duration, opt ...Option) Pipeable[T, []T] {
	o := mustOptions(opt...)

	return func(ctx context.Context, in Stream[T]) Stream[[]T] {
		out := make(chan Item[[]T], o.bufferSize)

		go func() {
			defer close(out)

			batch := make([]T, 0, n)

			copyBatch := func() []T {
				batchCopy := make([]T, len(batch))
				copy(batchCopy, batch)
				return batchCopy
			}

			for {
				select {
				case item, ok := <-in:
					if !ok {
						if len(batch) > 0 {
							select {
							case out <- Item[[]T]{Val: copyBatch()}:
							case <-ctx.Done():
								out <- Item[[]T]{Err: ctx.Err()}
							}
						}
						return
					}

					if item.Err != nil {
						select {
						case out <- Item[[]T]{Err: item.Err}:
						case <-ctx.Done():
							out <- Item[[]T]{Err: ctx.Err()}
							return
						}
						if o.stopOnError {
							select {
							case out <- Item[[]T]{Val: copyBatch()}:
								batch = batch[:0]
							case <-ctx.Done():
								out <- Item[[]T]{Err: ctx.Err()}
								return
							}
							return
						}
						continue
					}

					batch = append(batch, item.Val)

					if len(batch) == n {
						select {
						case out <- Item[[]T]{Val: copyBatch()}:
							batch = batch[:0]
						case <-ctx.Done():
							out <- Item[[]T]{Err: ctx.Err()}
							return
						}
					}

				case <-time.After(maxWait):
					if len(batch) > 0 {
						select {
						case out <- Item[[]T]{Val: copyBatch()}:
							batch = batch[:0]
						case <-ctx.Done():
							out <- Item[[]T]{Err: ctx.Err()}
							return
						}
					}
				case <-ctx.Done():
					out <- Item[[]T]{Err: ctx.Err()}
					return
				}
			}
		}()

		return out
	}
}
