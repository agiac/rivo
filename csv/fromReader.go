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
	return rivo.FromFunc(func(ctx context.Context, errs chan<- error) ([]string, bool, bool) {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			return []string{}, false, false
		}

		if err != nil {
			select {
			case <-ctx.Done():
				return nil, false, false
			case errs <- err:
				return nil, true, true
			}
		}

		return record, false, true
	})
}
