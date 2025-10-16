package io

import (
	"context"
	"io"

	"github.com/agiac/rivo"
)

// TODO: consider using ForEachOutput function

// ToWriter returns a pipeline that writes to an io.Writer.
func ToWriter(w io.Writer) rivo.Pipeline[[]byte, int] {
	return rivo.Map[[]byte, int](func(ctx context.Context, v []byte) (int, error) {
		return w.Write(v)
	})
}
