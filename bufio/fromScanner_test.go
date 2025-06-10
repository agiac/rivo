package bufio_test

import (
	"bufio"
	"context"
	"strings"
	"testing"

	. "github.com/agiac/rivo/bufio"
	. "github.com/agiac/rivo/core"
	"github.com/stretchr/testify/assert"
)

func TestFromScanner(t *testing.T) {
	t.Run("small buffer", func(t *testing.T) {
		ctx := context.Background()

		s := "1\n2\n3\n4\n5\n"
		r := strings.NewReader(s)
		scanner := bufio.NewScanner(r)

		g := FromScanner(scanner)

		got := Collect(g(ctx, nil))

		want := []Item[[]byte]{
			{Val: []byte("1")},
			{Val: []byte("2")},
			{Val: []byte("3")},
			{Val: []byte("4")},
			{Val: []byte("5")},
		}

		assert.Equal(t, want, got)
	})

	t.Run("large buffer", func(t *testing.T) {
		ctx := context.Background()

		s := strings.Repeat("Hello World\n", 1000)
		r := strings.NewReader(s)
		scanner := bufio.NewScanner(r)

		g := FromScanner(scanner)

		got := Collect(g(ctx, nil))

		want := make([]Item[[]byte], 1000)
		for i := range want {
			want[i] = Item[[]byte]{Val: []byte("Hello World")}
		}

		assert.Equal(t, want, got)
	})
}
