package bufio_test

import (
	"bufio"
	"bytes"
	"context"
	. "github.com/agiac/rivo"
	"testing"

	. "github.com/agiac/rivo/bufio"
	"github.com/stretchr/testify/assert"
)

func TestToWriter(t *testing.T) {
	t.Run("flush data", func(t *testing.T) {
		ctx := context.Background()

		in := Of([]byte("hello "), []byte("world"))

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		write := ToWriter(w)

		Collect(Pipe(in, write)(ctx, nil, nil))

		assert.Equal(t, "hello world", buf.String())
	})
}
