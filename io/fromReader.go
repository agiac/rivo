package io

import (
	"context"
	"io"

	"github.com/agiac/rivo"
)

// TODO: consider using ForEachOutput function

// FromReader returns a worker that reads from an io.Reader.
func FromReader(r io.Reader) rivo.Worker[rivo.None, []byte] {
	return func(ctx context.Context, _ rivo.Stream[rivo.None], errs chan<- error) rivo.Stream[[]byte] {
		out := make(chan []byte)

		go func() {
			defer close(out)

			buf := make([]byte, 1024)

			for {
				n, err := r.Read(buf)
				if err != nil {
					if err == io.EOF {
						return
					}
					select {
					case <-ctx.Done():
					case errs <- err:
					}
					continue
				}

				val := make([]byte, n)
				copy(val, buf[:n])

				select {
				case <-ctx.Done():
					return
				case out <- val:
				}
			}
		}()

		return out
	}
}
