package csv

import (
	"context"
	"encoding/csv"

	"github.com/agiac/rivo"
)

// ToWriter returns a pipeline that writes to a csv.Writer. Only errors from the
// csv.Writer are passed to the output stream.
func ToWriter(w *csv.Writer, opts ...rivo.Option) rivo.Pipeline[[]string, struct{}] {
	return func(ctx context.Context, in rivo.Stream[[]string]) rivo.Stream[struct{}] {
		out := make(chan rivo.Item[struct{}])

		go func() {
			defer close(out)
			defer w.Flush()

			for item := range rivo.OrDone(ctx, in) {
				if item.Err != nil {
					select {
					case <-ctx.Done():
						out <- rivo.Item[struct{}]{Err: ctx.Err()}
						return
					case out <- rivo.Item[struct{}]{Err: item.Err}:
					}
					continue
				}

				if err := w.Write(item.Val); err != nil {
					select {
					case <-ctx.Done():
						out <- rivo.Item[struct{}]{Err: ctx.Err()}
						return
					case out <- rivo.Item[struct{}]{Err: err}:

					}
				}
			}
		}()

		return out
	}
}
