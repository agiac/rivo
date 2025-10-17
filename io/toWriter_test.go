package io_test

import (
	"bytes"
	"context"
	. "github.com/agiac/rivo"
	"testing"

	. "github.com/agiac/rivo/io"
	"github.com/stretchr/testify/assert"
)

func TestToWriter(t *testing.T) {
	ctx := context.Background()

	in := Of([]byte("hello "), []byte("world"))

	var buf bytes.Buffer

	Collect(Pipe(in, ToWriter(&buf))(ctx, nil, nil))

	assert.Equal(t, "hello world", buf.String())
}
