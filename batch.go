package rivo

import (
	"context"
	"time"
)

// TODO: review this

func Batch[T any](batchSize int, maxWait time.Duration) Pipeable[T, []Item[T]] {
	return func(ctx context.Context, in Stream[T]) Stream[[]Item[T]] {
		out := make(chan Item[[]Item[T]])

		go func() {
			defer close(out)

			var batch []Item[T]

			sendBatch := func() {
				if len(batch) == 0 {
					return
				}

				out <- Item[[]Item[T]]{Val: batch}

				batch = nil
			}

			for {
				select {
				case <-ctx.Done():
					sendBatch()
					out <- Item[[]Item[T]]{Err: ctx.Err()}
					return
				case <-time.After(maxWait):
					sendBatch()
				case elem, ok := <-in:
					if !ok {
						sendBatch()
						return
					}

					batch = append(batch, elem)

					if len(batch) == batchSize {
						sendBatch()
					}
				}
			}
		}()

		return out
	}
}

func Batch2[T any](batchSize int, maxWait time.Duration, opt ...Option) Pipeable[T, []T] {
	o := mustOptions(opt...)

	return func(ctx context.Context, in Stream[T]) Stream[[]T] {
		out := make(chan Item[[]T], o.bufferSize)

		go func() {
			defer close(out)

			var batch []T

			sendBatch := func() {
				if len(batch) == 0 {
					return
				}

				out <- Item[[]T]{Val: batch}

				batch = nil
			}

			defer sendBatch()

			for item := range OrDone(ctx, in) {
				if item.Err != nil {
					select {
					case <-ctx.Done():
						out <- Item[[]T]{Err: ctx.Err()}
					case out <- Item[[]T]{Err: item.Err}:
						if o.stopOnError {
							return
						} else {
							continue
						}
					}
				}

				batch = append(batch, item.Val)

				select {
				case <-ctx.Done():
					out <- Item[[]T]{Err: ctx.Err()}
					return
				case <-time.After(maxWait):
					sendBatch()
				default:
				}

				if len(batch) == batchSize {
					sendBatch()
				}
			}
		}()

		return out
	}
}
