package rivo

import (
	"context"
)

// Flatten returns a Pipeable that flattens a Stream of slices into a Stream of individual items.
func Flatten[T any](opt ...Option) Pipeable[[]T, T] {
	o := mustOptions(opt...)

	return func(ctx context.Context, in Stream[[]T]) Stream[T] {
		out := make(chan Item[T], o.bufferSize)

		go func() {
			defer close(out)
			defer beforeClose(ctx, out, o)

			for {
				select {
				case item, ok := <-in:
					if !ok {
						return
					}

					if item.Err != nil {
						select {
						case out <- Item[T]{Err: item.Err}:
							if o.stopOnError {
								return
							} else {
								continue
							}
						case <-ctx.Done():
							out <- Item[T]{Err: ctx.Err()}
							return
						}
					}

					for _, val := range item.Val {
						select {
						case out <- Item[T]{Val: val}:
						case <-ctx.Done():
							out <- Item[T]{Err: ctx.Err()}
							return
						}
					}
				case <-ctx.Done():
					out <- Item[T]{Err: ctx.Err()}
					return
				}
			}
		}()

		return out
	}
}
