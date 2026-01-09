package io

import (
	"context"
	"io"

	"github.com/agiac/rivo"
)

// TODO: consider using ForEachOutput function

// ToWriter returns a worker that writes to an io.Writer.
func ToWriter(w io.Writer) rivo.Worker[[]byte, int] {
	return rivo.Map[[]byte, int](func(ctx context.Context, v []byte) (int, error) {
		return w.Write(v)
	})
}
