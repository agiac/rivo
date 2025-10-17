package bufio

import (
	"bufio"
	"context"
	"github.com/agiac/rivo"

	"github.com/agiac/rivo/io"
)

// TODO: consider using ForEachOutput function

// ToWriter returns a pipeline that writes to a bufio.Writer.
func ToWriter(w *bufio.Writer) rivo.Pipeline[[]byte, int] {
	return func(ctx context.Context, in rivo.Stream[[]byte], errs chan<- error) rivo.Stream[int] {
		out := make(chan int)

		go func() {
			defer close(out)
			defer func() {
				if err := w.Flush(); err != nil {
					select {
					case <-ctx.Done():
					case errs <- err:
					}
				}
			}()

			for item := range io.ToWriter(w)(ctx, in, errs) {
				select {
				case <-ctx.Done():
					return
				case out <- item:
				}
			}

		}()

		return out
	}
}
