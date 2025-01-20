package rivo

import (
	"context"
	"fmt"
	"time"
)

type BatchOptions struct {
	maxWait    time.Duration
	bufferSize int
}

type BatchOption func(*BatchOptions)

func BatchMaxWait(d time.Duration) BatchOption {
	return func(o *BatchOptions) {
		o.maxWait = d
	}
}

func BatchBufferSize(n int) BatchOption {
	return func(o *BatchOptions) {
		o.bufferSize = n
	}
}

var batchDefaultOptions = BatchOptions{
	maxWait:    1 * time.Second,
	bufferSize: 0,
}

func applyBatchOptions(opt []BatchOption) BatchOptions {
	opts := batchDefaultOptions
	for _, o := range opt {
		o(&opts)
	}
	return opts
}

func validateBatchOptions(opt BatchOptions) error {
	if opt.maxWait <= 0 {
		return fmt.Errorf("maxWait must be greater than 0")
	}

	if opt.bufferSize < 0 {
		return fmt.Errorf("bufferSize must be greater than or equal to 0")
	}

	return nil
}

// Batch returns a Pipeline that batches items from the input Stream into slices of n items.
// If the batch is not full after maxWait, it will be sent anyway.
// Any error in the input Stream will be propagated to the output Stream immediately.
func Batch[T any](n int, opt ...BatchOption) Pipeline[T, []T] {
	o := applyBatchOptions(opt)
	if err := validateBatchOptions(o); err != nil {
		panic(fmt.Errorf("invalid batch options: %v", err))
	}

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

						continue
					}

					batch = append(batch, item.Val)

					if len(batch) == n {
						if exit := sendBatch(); exit {
							return
						}
					}
				case <-time.After(o.maxWait):
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
