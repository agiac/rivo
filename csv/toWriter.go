package csv

import (
	"context"
	"encoding/csv"

	"github.com/agiac/rivo"
)

// ToWriter returns a pipeable that writes to a csv.Writer. Only errors from the
// csv.Writer are passed to the output stream.
func ToWriter(w *csv.Writer, opts ...rivo.Option) rivo.Pipeable[[]string, struct{}] {
	writeRow := func(ctx context.Context, i rivo.Item[[]string]) error {
		if i.Err != nil {
			return i.Err
		}

		return w.Write(i.Val)
	}

	flush := func(ctx context.Context) error {
		w.Flush()
		return w.Error()
	}

	return rivo.ForEach(writeRow, append(opts, rivo.WithOnBeforeClose(flush))...)
}
