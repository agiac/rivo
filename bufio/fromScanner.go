package bufio

import (
	"bufio"
	"context"

	"github.com/agiac/rivo"
)

// FromScanner returns a generator pipeable that reads from a bufio.Scanner.
// It's not thread-safe to use a pool size greater than 1.
func FromScanner(s *bufio.Scanner, opt ...rivo.Option) rivo.Pipeable[rivo.None, []byte] {
	return rivo.FromFunc[[]byte](func(ctx context.Context) ([]byte, error) {
		if !s.Scan() {
			if err := s.Err(); err != nil {
				return nil, err
			}

			return nil, rivo.ErrEOS
		}

		return s.Bytes(), nil
	}, opt...)
}
