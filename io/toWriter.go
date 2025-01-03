package io

import (
	"context"
	"io"

	"github.com/agiac/rivo"
)

// ToWriter returns a rivo.Transformer that writes to an io.Writer.
func ToWriter(w io.Writer, opt ...rivo.Option) rivo.Transformer[[]byte, int] {
	return rivo.Map(func(ctx context.Context, i rivo.Item[[]byte]) (int, error) {
		if i.Err != nil {
			return 0, i.Err
		}

		return w.Write(i.Val)
	}, opt...)
}
