package rivo_bufio_test

import (
	"bufio"
	"bytes"
	"context"
	. "github.com/agiac/rivo"
	. "github.com/agiac/rivo/rivo_bufio"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToWriter(t *testing.T) {
	t.Run("flush data", func(t *testing.T) {
		ctx := context.Background()

		in := Of([]byte("hello "), []byte("world"))

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		write := ToWriter(w)

		Collect(Pipe(in, write)(ctx, nil))

		assert.Equal(t, "hello world", buf.String())
	})

	t.Run("with on before close", func(t *testing.T) {
		ctx := context.Background()

		in := Of([]byte("hello "), []byte("world"))

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		beforeCloseCalled := false
		write := ToWriter(w, WithOnBeforeClose(func(ctx context.Context) error {
			beforeCloseCalled = true
			return nil
		}))

		Collect(Pipe(in, write)(ctx, nil))

		assert.Equal(t, "hello world", buf.String())
		assert.True(t, beforeCloseCalled)
	})
}
