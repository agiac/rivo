package core

import (
	"context"
	"fmt"
	"time"
)

// Batch returns a Pipeline that batches items from the input Stream into slices of n items.
// If the batch is not full after maxWait, it will be sent anyway.
// Any error in the input Stream will be propagated to the output Stream immediately.
func Batch[T any](n int, opt ...BatchOption) Pipeline[T, []T] {
	o := assertBatchOptions(opt)

	return func(ctx context.Context, in Stream[T]) Stream[[]T] {
		out := make(chan []T, o.bufferSize)

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
					case out <- copyBatch():
						batch = batch[:0]
					case <-ctx.Done():
						return true
					}
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

					batch = append(batch, item)

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
					return
				}
			}
		}()

		return out
	}
}

type batchOptions struct {
	maxWait    time.Duration
	bufferSize int
}

type BatchOption func(*batchOptions) error

func BatchMaxWait(d time.Duration) BatchOption {
	return func(o *batchOptions) error {
		if d <= 0 {
			return fmt.Errorf("maxWait must be greater than 0")
		}
		o.maxWait = d
		return nil
	}
}

func BatchBufferSize(n int) BatchOption {
	return func(o *batchOptions) error {
		if n < 0 {
			return fmt.Errorf("bufferSize must be greater than or equal to 0")
		}
		o.bufferSize = n
		return nil
	}
}

func newDefaultBatchOptions() *batchOptions {
	return &batchOptions{
		maxWait:    1 * time.Second,
		bufferSize: 0,
	}
}

func applyBatchOptions(opt []BatchOption) (*batchOptions, error) {
	opts := newDefaultBatchOptions()
	for _, o := range opt {
		if err := o(opts); err != nil {
			return opts, err
		}
	}
	return opts, nil
}

func assertBatchOptions(opt []BatchOption) *batchOptions {
	opts, err := applyBatchOptions(opt)
	if err != nil {
		panic(fmt.Errorf("invalid batch options: %v", err))
	}
	return opts
}
