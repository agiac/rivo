package bufio_test

import (
	"bufio"
	"context"
	. "github.com/agiac/rivo"
	. "github.com/agiac/rivo/bufio"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFromScanner(t *testing.T) {
	t.Run("small buffer", func(t *testing.T) {
		ctx := context.Background()

		s := "1\n2\n3\n4\n5\n"
		r := strings.NewReader(s)
		scanner := bufio.NewScanner(r)

		g := FromScanner(scanner)

		got := Collect(g(ctx, nil, nil))

		want := [][]byte{
			[]byte("1"),
			[]byte("2"),
			[]byte("3"),
			[]byte("4"),
			[]byte("5"),
		}

		assert.Equal(t, want, got)
	})

	t.Run("large buffer", func(t *testing.T) {
		ctx := context.Background()

		s := strings.Repeat("Hello World\n", 1000)
		r := strings.NewReader(s)
		scanner := bufio.NewScanner(r)

		g := FromScanner(scanner)

		got := Collect(g(ctx, nil, nil))

		want := make([][]byte, 1000)
		for i := range want {
			want[i] = []byte("Hello World")
		}

		assert.Equal(t, want, got)
	})
}
