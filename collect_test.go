package rivo_test

import (
	"context"
	"fmt"
	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleCollect() {
	ctx := context.Background()

	s := Of(1, 2, 3, 4, 5)(ctx, nil)

	for _, item := range Collect(s) {
		fmt.Println(item.Val)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func TestCollect(t *testing.T) {
	t.Run("collect till end of input stream", func(t *testing.T) {
		ctx := context.Background()

		in := Of(1, 2, 3, 4, 5)(ctx, nil)

		got := Collect(in)

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
		}

		assert.Equal(t, want, got)
	})
}

func TestCollectWithContext(t *testing.T) {
	t.Run("collect till end of input stream", func(t *testing.T) {
		ctx := context.Background()

		in := Of(1, 2, 3, 4, 5)(ctx, nil)

		got := CollectWithContext(ctx, in)

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
		}

		assert.Equal(t, want, got)
	})

	t.Run("collect till context is cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		g := Of(1, 2, 3, 4, 5)
		got := CollectWithContext(ctx, g(ctx, nil))

		assert.Lessf(t, len(got), 5, "expected less than 5 items due to context cancellation")
	})
}
