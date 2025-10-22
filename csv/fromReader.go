package csv

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/agiac/rivo"
)

// FromReader returns a generator pipeline that reads from a csv.Reader.
// It's not thread-safe to use a pool size greater than 1.
func FromReader(r *csv.Reader, opt ...FromReaderOption) rivo.Pipeline[rivo.None, []string] {
	o := assertFromReaderOptions(opt)

	return rivo.FromFunc(func(ctx context.Context) ([]string, bool, error) {
		// Discard the header line on first read if requested
		if o.discardHeader {
			o.discardHeader = false // Only discard once
			if _, err := r.Read(); err != nil && !errors.Is(err, io.EOF) {
				return nil, false, err
			}
		}

		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			return nil, false, nil
		}

		return record, true, err
	})
}

type fromReaderOptions struct {
	discardHeader bool
}

type FromReaderOption func(*fromReaderOptions) error

// DiscardHeader configures FromReader to skip the first line of the CSV file.
func DiscardHeader() FromReaderOption {
	return func(o *fromReaderOptions) error {
		o.discardHeader = true
		return nil
	}
}

func newDefaultFromReaderOptions() *fromReaderOptions {
	return &fromReaderOptions{
		discardHeader: false,
	}
}

func applyFromReaderOptions(opt []FromReaderOption) (*fromReaderOptions, error) {
	opts := newDefaultFromReaderOptions()
	for _, o := range opt {
		if err := o(opts); err != nil {
			return opts, err
		}
	}
	return opts, nil
}

func assertFromReaderOptions(opt []FromReaderOption) *fromReaderOptions {
	opts, err := applyFromReaderOptions(opt)
	if err != nil {
		panic(fmt.Errorf("invalid fromReader options: %v", err))
	}
	return opts
}
