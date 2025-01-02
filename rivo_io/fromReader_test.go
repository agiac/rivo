package rivo_io_test

import (
	"context"
	. "github.com/agiac/rivo/rivo_io"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFromReader(t *testing.T) {
	t.Run("small buffer", func(t *testing.T) {
		ctx := context.Background()

		g := FromReader(strings.NewReader("Hello World"))

		var got string
		for item := range g(ctx, nil) {
			if item.Err != nil {
				assert.Fail(t, "unexpected error", item.Err)
				return
			}
			got += string(item.Val)
		}

		assert.Equal(t, "Hello World", got)
	})

	t.Run("large buffer", func(t *testing.T) {
		ctx := context.Background()

		s := strings.Repeat("Hello World", 1000)
		g := FromReader(strings.NewReader(s))

		var got string
		for item := range g(ctx, nil) {
			if item.Err != nil {
				assert.Fail(t, "unexpected error", item.Err)
				return
			}
			got += string(item.Val)
		}

		assert.Equal(t, s, got)
	})

}
