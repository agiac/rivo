package bufio

import (
	"bufio"
	"context"
	"github.com/agiac/rivo"
)

// TODO: consider using ForEachOutput function

// FromScanner returns a generator pipeline that reads from a bufio.Scanner.
func FromScanner(s *bufio.Scanner) rivo.Pipeline[rivo.None, []byte] {
	return rivo.FromFunc[[]byte](func(ctx context.Context, errs chan<- error) ([]byte, bool, bool) {
		if !s.Scan() {
			if err := s.Err(); err != nil {
				select {
				case <-ctx.Done():
					return nil, false, false
				case errs <- err:
					return nil, true, true
				}
			}

			return nil, false, false // Stop the generator when no more data is available
		}

		return s.Bytes(), false, true // Return the scanned bytes
	})
}
