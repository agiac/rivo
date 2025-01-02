package rivo_io_test

import (
	"bytes"
	"context"
	. "github.com/agiac/rivo"
	. "github.com/agiac/rivo/rivo_io"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToWriter(t *testing.T) {
	ctx := context.Background()

	in := Of([]byte("hello "), []byte("world"))

	var buf bytes.Buffer

	Collect(Pipe(in, ToWriter(&buf))(ctx, nil))

	assert.Equal(t, "hello world", buf.String())
}
