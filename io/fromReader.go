package io

import (
	"context"
	"io"

	"github.com/agiac/rivo"
)

// FromReader returns a pipeline that reads from an io.Reader.
func FromReader(r io.Reader) rivo.Pipeline[rivo.None, []byte] {
	return func(ctx context.Context, _ rivo.Stream[rivo.None]) rivo.Stream[[]byte] {
		out := make(chan rivo.Item[[]byte])

		go func() {
			defer close(out)

			buf := make([]byte, 1024)

			for {
				n, err := r.Read(buf)
				if err != nil {
					if err == io.EOF {
						return
					}
					out <- rivo.Item[[]byte]{Err: err}
					continue
				}

				val := make([]byte, n)
				copy(val, buf[:n])

				select {
				case <-ctx.Done():
					out <- rivo.Item[[]byte]{Err: ctx.Err()}
					return
				case out <- rivo.Item[[]byte]{Val: val}:
				}
			}
		}()

		return out
	}
}
