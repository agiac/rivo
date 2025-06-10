package io

import (
	"context"
	"io"

	rivo "github.com/agiac/rivo/core"
)

// TODO: consider using ForEachOutput function

// ToWriter returns a pipeline that writes to an io.Writer.
func ToWriter(w io.Writer) rivo.Pipeline[[]byte, rivo.Item[int]] {
	return rivo.Map[[]byte, rivo.Item[int]](func(ctx context.Context, v []byte) rivo.Item[int] {
		n, err := w.Write(v)
		return rivo.Item[int]{Val: n, Err: err}
	})
}
