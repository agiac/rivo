package rivo

import (
	"context"
	"time"
)

// Batch returns a Pipeline that batches items from the input Stream into slices of n items.
// If the batch is not full after maxWait, it will be sent anyway.
// Any error in the input Stream will be propagated to the output Stream immediately.
func Batch[T any](n int, maxWait time.Duration, opt ...Option) Pipeline[T, []T] {
	o := mustOptions(opt...)

	return func(ctx context.Context, in Stream[T]) Stream[[]T] {
		out := make(chan Item[[]T], o.bufferSize)

		go func() {
			defer close(out)
			defer beforeClose(ctx, out, o)

			batch := make([]T, 0, n)

			copyBatch := func() []T {
				batchCopy := make([]T, len(batch))
				copy(batchCopy, batch)
				return batchCopy
			}

			sendBatch := func() (exit bool) {
				if len(batch) > 0 {
					select {
					case out <- Item[[]T]{Val: copyBatch()}:
						batch = batch[:0]
					case <-ctx.Done():
						out <- Item[[]T]{Err: ctx.Err()}
						return true
					}
				}
				return false
			}

			sendError := func(err error) (exit bool) {
				select {
				case out <- Item[[]T]{Err: err}:
				case <-ctx.Done():
					out <- Item[[]T]{Err: ctx.Err()}
					return true
				}
				return false
			}

			for {
				select {
				case item, ok := <-in:
					if !ok {
						sendBatch()
						return
					}

					if item.Err != nil {
						if exit := sendError(item.Err); exit {
							return
						}

						if o.stopOnError {
							sendBatch()
							return
						}

						continue
					}

					batch = append(batch, item.Val)

					if len(batch) == n {
						if exit := sendBatch(); exit {
							return
						}
					}
				case <-time.After(maxWait):
					if exit := sendBatch(); exit {
						return
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
