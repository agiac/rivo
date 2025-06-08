package rivo_test

import (
	"context"
	"fmt"
	"github.com/agiac/rivo/core"
	"testing"
	"time"

	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

func ExampleBatch() {
	ctx := context.Background()

	in := Of(1, 2, 3, 4, 5)

	b := Batch[int](2)

	p := core.Pipe(in, b)

	for item := range p(ctx, nil) {
		fmt.Printf("%v\n", item.Val)
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

		got := core.Collect(core.Pipe(in, b)(ctx, nil))

		want := []Item[[]int]{
			{Val: []int{1, 2}},
			{Val: []int{3, 4}},
			{Val: []int{5, 6}},
		}

		assert.Equal(t, want, got)
	})

	t.Run("batch items by time", func(t *testing.T) {
		ctx := context.Background()

		in := make(chan Item[int])

		go func() {
			defer close(in)
			in <- Item[int]{Val: 1}
			time.Sleep(200 * time.Millisecond)
			in <- Item[int]{Val: 2}
			time.Sleep(200 * time.Millisecond)
			in <- Item[int]{Val: 3}
		}()

		b := Batch[int](10, BatchMaxWait(100*time.Millisecond))

		got := core.Collect(b(ctx, in))

		want := []Item[[]int]{
			{Val: []int{1}},
			{Val: []int{2}},
			{Val: []int{3}},
		}

		assert.Equal(t, want, got)
	})

	t.Run("don't batch errors", func(t *testing.T) {
		ctx := context.Background()

		in := make(chan Item[int])

		go func() {
			defer close(in)
			in <- Item[int]{Val: 1}
			in <- Item[int]{Err: fmt.Errorf("error")}
			in <- Item[int]{Val: 2}
		}()

		b := Batch[int](2)

		got := core.Collect(b(ctx, in))

		want := []Item[[]int]{
			{Err: fmt.Errorf("error")},
			{Val: []int{1, 2}},
		}

		assert.Equal(t, want, got)
	})

	t.Run("context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		g := Of(1, 2, 3, 4, 5)
		b := Batch[int](2)

		got := core.Collect(core.Pipe(g, b)(ctx, nil))

		assert.Lessf(t, len(got), 3, "expected less than 3 items due to context cancellation")
	})

	t.Run("with buffer size", func(t *testing.T) {
		ctx := context.Background()

		in := Of(1, 2, 3)

		b := Batch[int](2, BatchBufferSize(3))

		out := core.Pipe(in, b)(ctx, nil)

		got := core.Collect(out)

		want := []Item[[]int]{
			{Val: []int{1, 2}},
			{Val: []int{3}},
		}

		assert.Equal(t, 3, cap(out))
		assert.Equal(t, want, got)
	})
}
