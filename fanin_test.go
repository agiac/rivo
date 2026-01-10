package rivo

import (
	"context"
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleFanin() {
	ctx := context.Background()

	g1 := Of("Hello", "World")
	g2 := Of("Foo", "Bar")

	gg := Fanin(g1, g2)

	res := Collect(gg(ctx, nil, nil))

	slices.Sort(res)

	fmt.Println(res)

	// Output:
	// [Bar Foo Hello World]
}

func TestFanin(t *testing.T) {
	t.Run("fanin combines outputs", func(t *testing.T) {
		ctx := context.Background()

		g1 := Of(1, 2)
		g2 := Of(3, 4)
		g3 := Of(5)

		out := Fanin(g1, g2, g3)(ctx, nil, nil)
		got := Collect(out)

		assert.ElementsMatch(t, []int{1, 2, 3, 4, 5}, got)
	})

	t.Run("fanin with no generators returns closed channel", func(t *testing.T) {
		ctx := context.Background()

		s := Collect(Fanin[int]()(ctx, nil, nil))
		assert.Empty(t, s)
	})

	t.Run("fanin with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		g1 := Of(1, 2, 3, 4, 5)
		g2 := Of(6, 7, 8, 9, 10)

		out := Fanin(g1, g2)(ctx, nil, nil)
		got := Collect(out)

		assert.Less(t, len(got), 5, "expected few or no items due to context cancellation")
	})
}
