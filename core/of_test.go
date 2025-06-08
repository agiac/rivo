package core_test

import (
	"context"
	"fmt"
	"testing"

	. "github.com/agiac/rivo/core"
	"github.com/stretchr/testify/assert"
)

func ExampleOf() {
	ctx := context.Background()

	in := Of(1, 2, 3, 4, 5)

	s := in(ctx, nil)

	for item := range s {
		fmt.Println(item)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func TestOf(t *testing.T) {
	t.Run("create stream from items", func(t *testing.T) {
		ctx := context.Background()

		p := Of(1, 2, 3, 4, 5)

		got := Collect(p(ctx, nil))

		want := []int{1, 2, 3, 4, 5}

		assert.Equal(t, want, got)
	})

	t.Run("context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		in := Of(1, 2, 3, 4, 5)(ctx, nil)

		got := Collect(in)

		assert.Lessf(t, len(got), 5, "should not collect all items when context is cancelled, got: %v", got)
	})
}
