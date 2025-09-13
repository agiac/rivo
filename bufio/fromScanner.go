package bufio

import (
	"bufio"
	"context"
	"fmt"
	"github.com/agiac/rivo"
)

// TODO: consider using ForEachOutput function

// FromScanner returns a generator pipeline that reads from a bufio.Scanner.
func FromScanner(s *bufio.Scanner) rivo.Pipeline[rivo.None, []byte] {
	return rivo.Pipeline[rivo.None, []byte](rivo.FromFunc[[]byte](func(ctx context.Context) ([]byte, bool) {
		if !s.Scan() {
			if err := s.Err(); err != nil {
				// TODO: handle
				fmt.Printf("FromScanner: scanner error: %v\n", err)
				return nil, true
			}

			return nil, false // Stop the generator when no more data is available
		}

		return s.Bytes(), true // Return the scanned bytes
	}))
}
