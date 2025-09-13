package rivo_test

import (
	"context"
	"fmt"
	. "github.com/agiac/rivo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExamplePipe() {
	ctx := context.Background()

	a := Of(1, 2, 3, 4, 5)

	p := Pipe(a, Map(func(ctx context.Context, n int) int {
		return n + 1
	}))

	s := p(ctx, nil, nil)

	for item := range s {
		fmt.Println(item)
	}

	// Output:
	// 2
	// 3
	// 4
	// 5
	// 6
}

func TestPipes(t *testing.T) {
	var addOne = Map(func(ctx context.Context, n int) int {
		return n + 1
	})

	t.Run("pipe", func(t *testing.T) {
		t.Run("pipe all values", func(t *testing.T) {
			ctx := context.Background()

			a := Of(1, 2, 3, 4, 5)

			p := Pipe(a, addOne)

			got := Collect(p(ctx, nil, nil))

			want := []int{2, 3, 4, 5, 6}

			assert.Equal(t, want, got)
		})

		t.Run("context cancelled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			a := Of(1, 2, 3, 4, 5)

			p := Pipe(a, addOne)

			got := Collect(p(ctx, nil, nil))

			assert.Lessf(t, len(got), 5, "should not collect all items when context is cancelled, got: %v", got)
		})
	})

	t.Run("pipe2", func(t *testing.T) {
		t.Run("pipe all values", func(t *testing.T) {
			ctx := context.Background()

			a := Of(1, 2, 3, 4, 5)

			p := Pipe2(a, addOne)

			got := Collect(p(ctx, nil, nil))

			want := []int{
				2, 3, 4, 5, 6}

			assert.Equal(t, want, got)
		})

		t.Run("context cancelled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			a := Of(1, 2, 3, 4, 5)

			p := Pipe2(a, addOne)

			got := Collect(p(ctx, nil, nil))

			assert.Lessf(t, len(got), 5, "should not collect all items when context is cancelled, got: %v", got)
		})
	})

	t.Run("pipe3", func(t *testing.T) {
		t.Run("pipe all values", func(t *testing.T) {
			ctx := context.Background()

			a := Of(1, 2, 3, 4, 5)

			p := Pipe3(a, addOne, addOne)

			got := Collect(p(ctx, nil, nil))

			want := []int{
				3, 4, 5, 6, 7}

			assert.Equal(t, want, got)
		})

		t.Run("context cancelled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			a := Of(1, 2, 3, 4, 5)

			p := Pipe3(a, addOne, addOne)

			got := Collect(p(ctx, nil, nil))

			assert.Lessf(t, len(got), 5, "should not collect all items when context is cancelled, got: %v", got)
		})
	})

	t.Run("pipe4", func(t *testing.T) {
		t.Run("pipe all values", func(t *testing.T) {
			ctx := context.Background()

			a := Of(1, 2, 3, 4, 5)

			p := Pipe4(a, addOne, addOne, addOne)

			got := Collect(p(ctx, nil, nil))

			want := []int{
				4, 5, 6, 7, 8}

			assert.Equal(t, want, got)
		})

		t.Run("context cancelled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			a := Of(1, 2, 3, 4, 5)

			p := Pipe4(a, addOne, addOne, addOne)

			got := Collect(p(ctx, nil, nil))

			assert.Lessf(t, len(got), 5, "should not collect all items when context is cancelled, got: %v", got)
		})
	})

	t.Run("pipe5", func(t *testing.T) {
		t.Run("pipe all values", func(t *testing.T) {
			ctx := context.Background()

			a := Of(1, 2, 3, 4, 5)

			p := Pipe5(a, addOne, addOne, addOne, addOne)

			got := Collect(p(ctx, nil, nil))

			want := []int{
				5, 6, 7, 8, 9}

			assert.Equal(t, want, got)
		})

		t.Run("context cancelled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			a := Of(1, 2, 3, 4, 5)

			p := Pipe5(a, addOne, addOne, addOne, addOne)

			got := Collect(p(ctx, nil, nil))

			assert.Lessf(t, len(got), 5, "should not collect all items when context is cancelled, got: %v", got)
		})
	})
}
