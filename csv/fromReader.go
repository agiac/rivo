package csv

import (
	"context"
	"encoding/csv"
	"errors"
	"github.com/agiac/rivo"
	"io"
)

// TODO: Add support to discard the header line.

// FromReader returns a generator pipeline that reads from a csv.Reader.
// It's not thread-safe to use a pool size greater than 1.
func FromReader(r *csv.Reader) rivo.Pipeline[rivo.None, []string] {
	return rivo.Pipeline[rivo.None, []string](rivo.FromFunc(func(ctx context.Context) ([]string, bool) {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			return []string{}, false
		}

		if err != nil { // TODO: handle error
			return nil, true
		}

		return record, true
	}))
}
