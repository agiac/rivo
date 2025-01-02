package rivo_io

import (
	"context"
	"github.com/agiac/rivo"
	"io"
)

// ToWriter returns a pipeable that writes to an io.Writer.
func ToWriter(w io.Writer, opt ...rivo.Option) rivo.Pipeable[[]byte, int] {
	return rivo.Map(func(ctx context.Context, i rivo.Item[[]byte]) (int, error) {
		if i.Err != nil {
			return 0, i.Err
		}

		return w.Write(i.Val)
	}, opt...)
}
