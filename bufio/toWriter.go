package bufio

import (
	"bufio"
	"context"

	"github.com/agiac/rivo"
	"github.com/agiac/rivo/io"
)

// ToWriter returns a rivo.Transformer that writes to a bufio.Writer.
// It's not thread-safe to use a pool size greater than 1.
func ToWriter(w *bufio.Writer, opt ...rivo.Option) rivo.Transformer[[]byte, int] {
	return io.ToWriter(w, append(opt, rivo.WithOnBeforeClose(func(ctx context.Context) error {
		return w.Flush()
	}))...)
}
