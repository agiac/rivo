package core_test

import (
	"context"
	"fmt"
	. "github.com/agiac/rivo/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleMap() {
	ctx := context.Background()

	in := Of(1, 2, 3, 4, 5)

	double := Map(func(ctx context.Context, n int) int {
		return n * 2
	})

	p := Pipe(in, double)

	s := p(ctx, nil)

	for n := range s {
		fmt.Println(n)
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

		mapFn := func(ctx context.Context, n int) int {
			return n + 1
		}

		g := Of(1, 2, 3, 4, 5)

		m := Map(mapFn)

		got := Collect(Pipe(g, m)(ctx, nil))

		want := []int{2, 3, 4, 5, 6}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx := context.Background()

		mapFn := func(ctx context.Context, n int) int {
			return n + 1
		}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		g := Of(1, 2, 3, 4, 6)
		m := Map(mapFn)

		got := Collect(Pipe(g, m)(ctx, nil))

		assert.Lessf(t, len(got), 3, "expected less than 3 items due to context cancellation")
	})

	t.Run("with buffer size", func(t *testing.T) {
		ctx := context.Background()

		mapFn := func(ctx context.Context, n int) int {
			return n + 1
		}

		in := make(chan int)

		go func() {
			defer close(in)
			in <- 1
			in <- 2
			in <- 3
		}()

		m := Map(mapFn, MapBufferSize(3))

		out := m(ctx, in)

		got := Collect(out)

		want := []int{2, 3, 4}

		assert.Equal(t, 3, cap(out))
		assert.Equal(t, want, got)
	})

	t.Run("with pool size", func(t *testing.T) {
		ctx := context.Background()

		mapFn := func(ctx context.Context, n int) int {
			return n + 1
		}

		in := Of(1, 2, 3, 4, 5)

		m := Map(mapFn, MapPoolSize(3))

		got := Collect(Pipe(in, m)(ctx, nil))

		want := []int{2, 3, 4, 5, 6}

		assert.ElementsMatch(t, want, got)
	})
}
