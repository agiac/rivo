package bufio

import (
	"bufio"
	"context"
	"github.com/agiac/rivo"
)

// TODO: consider using ForEachOutput function

// FromScanner returns a generator pipeline that reads from a bufio.Scanner.
func FromScanner(s *bufio.Scanner) rivo.Pipeline[rivo.None, rivo.Item[[]byte]] {
	return rivo.Pipeline[rivo.None, rivo.Item[[]byte]](rivo.FromFunc[rivo.Item[[]byte]](func(ctx context.Context) (rivo.Item[[]byte], bool) {
		if !s.Scan() {
			if err := s.Err(); err != nil {
				return rivo.Item[[]byte]{Err: err}, true // Stop the generator on error
			}

			return rivo.Item[[]byte]{Val: nil}, false // Stop the generator when no more data is available
		}

		return rivo.Item[[]byte]{Val: s.Bytes()}, true // Return the scanned bytes
	}))
}
