package csv

import (
	"context"
	"encoding/csv"

	"github.com/agiac/rivo"
)

// ToWriter returns a sync pipeline that writes to a csv.Writer. Each input []string is written as a row in the CSV.
// Any errors encountered during writing are sent to the errs channel.
func ToWriter(w *csv.Writer) rivo.Sync[[]string] {
	return rivo.ForEachOutput(
		func(ctx context.Context, val []string, out chan<- rivo.None, errs chan<- error) {
			if err := w.Write(val); err != nil {
				select {
				case <-ctx.Done():
					return
				case errs <- err:
				}
			}
		},
		rivo.ForEachOutputOnBeforeClose(func(ctx context.Context) {
			w.Flush()
		}),
	)
}
