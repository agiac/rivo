package csv

import (
	"context"
	"encoding/csv"

	"github.com/agiac/rivo"
)

// ToWriter returns a pipeline that writes to a csv.Writer. Only errors from the
// csv.Writer are passed to the output stream.
func ToWriter(w *csv.Writer) rivo.Pipeline[[]string, error] {
	return rivo.ForEachOutput(
		func(ctx context.Context, val []string, out chan<- error, errs chan<- error) {
			if err := w.Write(val); err != nil {
				select {
				case <-ctx.Done():
					return
				case out <- err:
				}
			}
		},
		rivo.ForEachOutputOnBeforeClose(func(ctx context.Context) {
			w.Flush()
		}),
	)
}
