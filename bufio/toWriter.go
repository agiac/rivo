package bufio

import (
	"bufio"
	"context"

	"github.com/agiac/rivo"
	"github.com/agiac/rivo/io"
)

// ToWriter returns a pipeline that writes to a bufio.Writer.
func ToWriter(w *bufio.Writer) rivo.Pipeline[[]byte, int] {
	return func(ctx context.Context, in rivo.Stream[[]byte]) rivo.Stream[int] {
		out := make(chan rivo.Item[int])

		go func() {
			defer close(out)
			defer func() {
				if err := w.Flush(); err != nil {
					out <- rivo.Item[int]{Err: err}
				}
			}()

			for item := range io.ToWriter(w)(ctx, in) {
				select {
				case <-ctx.Done():
					out <- rivo.Item[int]{Err: ctx.Err()}
					return
				case out <- item:
				}
			}

		}()

		return out
	}
}
