package csv

import (
	"context"
	"encoding/csv"
	"errors"
	"io"

	"github.com/agiac/rivo"
)

// FromReader returns a generator pipeline that reads from a csv.Reader.
// It's not thread-safe to use a pool size greater than 1.
func FromReader(r *csv.Reader) rivo.Pipeline[rivo.None, []string] {
	return rivo.FromFunc(func(ctx context.Context) ([]string, error) {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			return nil, rivo.ErrEOS
		}

		if err != nil {
			return nil, err
		}

		return record, nil
	})
}
