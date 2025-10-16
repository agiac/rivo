package bufio

import (
	"bufio"
	"context"

	"github.com/agiac/rivo"
)

// TODO: consider using ForEachOutput function

// FromScanner returns a generator pipeline that reads from a bufio.Scanner.
func FromScanner(s *bufio.Scanner) rivo.Pipeline[rivo.None, []byte] {
	return rivo.FromFunc[[]byte](func(ctx context.Context) ([]byte, bool, error) {
		if !s.Scan() {
			return nil, false, s.Err() // Return any scanning error
		}

		return s.Bytes(), true, nil // Return the scanned bytes
	})
}
