package rivo_test

import (
	"context"
	"fmt"
	"testing"

	. "github.com/agiac/rivo"

	"github.com/stretchr/testify/assert"
)

func ExampleMap() {
	ctx := context.Background()

	in := Of(1, 2, 3, 4, 5)

	double := Map(func(ctx context.Context, n int) (int, error) {
		return n * 2, nil
	})

	p := Pipe(in, double)

	s := p(ctx, nil, nil)

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

		mapFn := func(ctx context.Context, n int) (int, error) {
			return n + 1, nil
		}

		g := Of(1, 2, 3, 4, 5)

		m := Map(mapFn)

		got := Collect(Pipe(g, m)(ctx, nil, nil))

		want := []int{2, 3, 4, 5, 6}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx := context.Background()

		mapFn := func(ctx context.Context, n int) (int, error) {
			return n + 1, nil
		}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		g := Of(1, 2, 3, 4, 6)
		m := Map(mapFn)

		got := Collect(Pipe(g, m)(ctx, nil, nil))

		assert.Lessf(t, len(got), 3, "expected less than 3 items due to context cancellation")
	})

	t.Run("with buffer size", func(t *testing.T) {
		ctx := context.Background()

		mapFn := func(ctx context.Context, n int) (int, error) {
			return n + 1, nil
		}

		in := make(chan int)

		go func() {
			defer close(in)
			in <- 1
			in <- 2
			in <- 3
		}()

		m := Map(mapFn, MapBufferSize(3))

		out := m(ctx, in, nil)

		got := Collect(out)

		want := []int{2, 3, 4}

		assert.Equal(t, 3, cap(out))
		assert.Equal(t, want, got)
	})

	t.Run("with pool size", func(t *testing.T) {
		ctx := context.Background()

		mapFn := func(ctx context.Context, n int) (int, error) {
			return n + 1, nil
		}

		in := Of(1, 2, 3, 4, 5)

		m := Map(mapFn, MapPoolSize(3))

		got := Collect(Pipe(in, m)(ctx, nil, nil))

		want := []int{2, 3, 4, 5, 6}

		assert.ElementsMatch(t, want, got)
	})
}
