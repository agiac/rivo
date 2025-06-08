package rivo_test

import (
	"context"
	"fmt"
	"github.com/agiac/rivo/core"
	"testing"

	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

func ExampleMap() {
	ctx := context.Background()

	in := Of(1, 2, 3, 4, 5)

	double := Map(func(ctx context.Context, i Item[int]) (int, error) {
		// Always check for errors
		if i.Err != nil {
			return 0, i.Err // Propagate the error
		}

		return i.Val * 2, nil
	})

	p := core.Pipe(in, double)

	s := p(ctx, nil)

	for item := range s {
		fmt.Println(item.Val)
	}

	// Output:
	// 2
	// 4
	// 6
	// 8
	// 10
}

func TestMap(t *testing.T) {
	t.Run("map all items", func(t *testing.T) {
		ctx := context.Background()

		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		}

		g := Of(1, 2, 3, 4, 5)

		m := Map(mapFn)

		got := core.Collect(core.Pipe(g, m)(ctx, nil))

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
			{Val: 6},
		}

		assert.Equal(t, want, got)
	})

	t.Run("map all items with error", func(t *testing.T) {
		ctx := context.Background()

		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			if i.Val == 3 {
				return 0, assert.AnError
			}
			return i.Val + 1, nil
		}

		g := Of(1, 2, 3, 4, 5)

		m := Map(mapFn)

		got := core.Collect(core.Pipe(g, m)(ctx, nil))

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Err: assert.AnError},
			{Val: 5},
			{Val: 6},
		}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx := context.Background()

		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			if i.Err != nil {
				return 0, i.Err
			}

			return i.Val + 1, nil
		}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		g := Of(1, 2, 3, 4, 6)
		m := Map(mapFn)

		got := core.Collect(core.Pipe(g, m)(ctx, nil))

		assert.Lessf(t, len(got), 3, "expected less than 3 items due to context cancellation")
	})

	t.Run("with buffer size", func(t *testing.T) {
		ctx := context.Background()

		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		}

		in := make(chan Item[int])

		go func() {
			defer close(in)
			in <- Item[int]{Val: 1}
			in <- Item[int]{Val: 2}
			in <- Item[int]{Val: 3}
		}()

		m := Map(mapFn, core.MapBufferSize(3))

		out := m(ctx, in)

		got := core.Collect(out)

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Val: 4},
		}

		assert.Equal(t, 3, cap(out))
		assert.Equal(t, want, got)
	})

	t.Run("with pool size", func(t *testing.T) {
		ctx := context.Background()

		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		}

		in := Of(1, 2, 3, 4, 5)

		m := Map(mapFn, core.MapPoolSize(3))

		got := core.Collect(core.Pipe(in, m)(ctx, nil))

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
			{Val: 6},
		}

		assert.ElementsMatch(t, want, got)
	})
}
