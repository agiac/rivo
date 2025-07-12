package rivo_test

import (
	"context"
	"fmt"
	. "github.com/agiac/rivo"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleBatch() {
	ctx := context.Background()

	in := Of(1, 2, 3, 4, 5)

	b := Batch[int](2)

	p := Pipe(in, b)

	for item := range p(ctx, nil) {
		fmt.Printf("%v\n", item)
	}

	// Output:
	// [1 2]
	// [3 4]
	// [5]
}

func TestBatch(t *testing.T) {
	t.Run("batch items by count", func(t *testing.T) {
		ctx := context.Background()

		in := Of(1, 2, 3, 4, 5, 6)

		b := Batch[int](2)

		got := Collect(Pipe(in, b)(ctx, nil))

		want := [][]int{
			{1, 2},
			{3, 4},
			{5, 6},
		}

		assert.Equal(t, want, got)
	})

	t.Run("batch items by time", func(t *testing.T) {
		ctx := context.Background()

		in := make(chan int)

		go func() {
			defer close(in)
			in <- 1
			time.Sleep(200 * time.Millisecond)
			in <- 2
			time.Sleep(200 * time.Millisecond)
			in <- 3
		}()

		b := Batch[int](10, BatchMaxWait(100*time.Millisecond))

		got := Collect(b(ctx, in))

		want := [][]int{{1}, {2}, {3}}

		assert.Equal(t, want, got)
	})

	t.Run("context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		in := Of(1, 2, 3, 4, 5)
		b := Batch[int](2)

		got := Collect(Pipe(in, b)(ctx, nil))

		assert.Lessf(t, len(got), 2, "should not batch items when context is cancelled")
	})

	t.Run("with buffer size", func(t *testing.T) {
		ctx := context.Background()

		in := Of(1, 2, 3)

		b := Batch[int](2, BatchBufferSize(3))

		out := Pipe(in, b)(ctx, nil)

		got := Collect(out)

		want := [][]int{{1, 2}, {3}}

		assert.Equal(t, 3, cap(out))
		assert.Equal(t, want, got)
	})
}
