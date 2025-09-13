package io

import (
	"context"
	"github.com/agiac/rivo"
	"io"
)

// TODO: consider using ForEachOutput function

// ToWriter returns a pipeline that writes to an io.Writer.
func ToWriter(w io.Writer) rivo.Pipeline[[]byte, int] {
	return rivo.Map[[]byte, int](func(ctx context.Context, v []byte) int {
		n, _ := w.Write(v) // TODO: handle error
		return n
	})
}
