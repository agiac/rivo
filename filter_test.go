package rivo_test

import (
	"context"
	"fmt"
	"testing"

	. "github.com/agiac/rivo"

	"github.com/stretchr/testify/assert"
)

func ExampleFilter() {
	ctx := context.Background()

	in := Of(1, 2, 3, 4, 5)

	even := Filter(func(ctx context.Context, n int) (bool, error) {
		return n%2 == 0, nil
	})

	p := Pipe(in, even)

	s := p(ctx, nil, nil)

	for item := range s {
		fmt.Println(item)
	}

	// Output:
	// 2
	// 4
}

func TestFilter(t *testing.T) {
	even := func(ctx context.Context, i int) (bool, error) {
		return i%2 == 0, nil
	}

	t.Run("filter all items", func(t *testing.T) {
		ctx := context.Background()

		g := Of(1, 2, 3, 4, 5)
		f := Filter(even)

		got := Collect(Pipe(g, f)(ctx, nil, nil))
		want := []int{2, 4}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		g := Of(1, 2, 3, 4, 5)
		f := Filter(even)

		got := Collect(f(ctx, g(ctx, nil, nil), nil))

		assert.Lessf(t, len(got), 5, "expected less than 5 items, got %d", len(got))
	})

	t.Run("with buffer size", func(t *testing.T) {
		ctx := context.Background()

		g := Of(1, 2, 3, 4, 5)
		f := Filter(even, FilterBufferSize(3))

		got := Collect(f(ctx, g(ctx, nil, nil), nil))
		want := []int{2, 4}

		assert.Equal(t, want, got)
	})

	t.Run("with pool size", func(t *testing.T) {
		ctx := context.Background()

		g := Of(1, 2, 3, 4, 5)
		f := Filter(even, FilterPoolSize(3))

		got := Collect(Pipe(g, f)(ctx, nil, nil))
		want := []int{2, 4}

		assert.ElementsMatch(t, want, got)
	})
}
