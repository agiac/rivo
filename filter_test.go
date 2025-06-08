package rivo_test

import (
	"context"
	"fmt"
	"github.com/agiac/rivo/core"
	"testing"

	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

func ExampleFilter() {
	ctx := context.Background()

	in := Of(1, 2, 3, 4, 5)

	onlyEven := Filter(func(ctx context.Context, i int) bool {
		return i%2 == 0
	})

	p := core.Pipe(in, onlyEven)

	s := p(ctx, nil)

	for item := range s {
		fmt.Println(item.Val)
	}

	// Output:
	// 2
	// 4
}

func TestFilter(t *testing.T) {
	t.Run("filter all items", func(t *testing.T) {
		ctx := context.Background()

		filterFn := func(ctx context.Context, i int) bool {
			return i%2 == 0
		}

		g := Of(1, 2, 3, 4, 5)

		f := Filter(filterFn)

		got := core.Collect(core.Pipe(g, f)(ctx, nil))

		want := []Item[int]{
			{Val: 2},
			{Val: 4},
		}

		assert.Equal(t, want, got)
	})

	t.Run("filter with error", func(t *testing.T) {
		ctx := context.Background()

		filterFn := func(ctx context.Context, i int) bool {
			return i%2 == 0
		}

		in := make(chan Item[int])
		go func() {
			defer close(in)
			in <- Item[int]{Val: 1}
			in <- Item[int]{Val: 2}
			in <- Item[int]{Err: assert.AnError} // Simulating an error
			in <- Item[int]{Val: 3}
			in <- Item[int]{Val: 4}
			in <- Item[int]{Val: 5}
		}()

		f := Filter(filterFn)

		got := core.Collect(f(ctx, in))

		want := []Item[int]{
			{Val: 2},
			{Err: assert.AnError},
			{Val: 4},
		}

		assert.Equal(t, want, got)
	})
}
