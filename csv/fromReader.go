package csv

import (
	"context"
	"encoding/csv"
	"errors"
	"io"

	rivo "github.com/agiac/rivo/core"
)

// FromReader returns a generator pipeline that reads from a csv.Reader.
// It's not thread-safe to use a pool size greater than 1.
func FromReader(r *csv.Reader) rivo.Pipeline[rivo.None, rivo.Item[[]string]] {
	return rivo.FromFunc(func(ctx context.Context) (rivo.Item[[]string], bool) {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			return rivo.Item[[]string]{}, false
		}

		if err != nil {
			return rivo.Item[[]string]{Err: err}, true
		}

		return rivo.Item[[]string]{Val: record}, true
	})
}
