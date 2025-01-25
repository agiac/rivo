package bufio

import (
	"bufio"
	"context"

	"github.com/agiac/rivo"
)

// FromScanner returns a generator pipeline that reads from a bufio.Scanner.
func FromScanner(s *bufio.Scanner) rivo.Pipeline[rivo.None, []byte] {
	return rivo.FromFunc[[]byte](func(ctx context.Context) ([]byte, error) {
		if !s.Scan() {
			if err := s.Err(); err != nil {
				return nil, err
			}

			return nil, rivo.ErrEOS
		}

		return s.Bytes(), nil
	})
}
