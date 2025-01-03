package io

import (
	"context"
	"io"

	"github.com/agiac/rivo"
)

// FromReader returns a pipeable that reads from an io.Reader.
// It's not thread-safe to use a pool size greater than 1.
func FromReader(r io.Reader, opt ...rivo.Option) rivo.Pipeable[struct{}, []byte] {
	buf := make([]byte, 1024)
	return rivo.FromFunc[[]byte](func(ctx context.Context) ([]byte, error) {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil, rivo.ErrEOS
			}
			return nil, err
		}

		val := make([]byte, n)
		copy(val, buf[:n])

		return val, nil
	}, opt...)
}
