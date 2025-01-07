package rivo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

func ExampleBatch() {
	ctx := context.Background()

	in := Of(1, 2, 3, 4, 5)

	b := Batch[int](2, time.Second)

	p := Pipe(in, b)

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

		b := Batch[int](2, time.Second)

		got := Collect(Pipe(in, b)(ctx, nil))

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

		b := Batch[int](10, 100*time.Millisecond)

		got := Collect(b(ctx, in))

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

		b := Batch[int](2, time.Second)

		got := Collect(b(ctx, in))

		want := []Item[[]int]{
			{Err: fmt.Errorf("error")},
			{Val: []int{1, 2}},
		}

		assert.Equal(t, want, got)
	})

	t.Run("context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		in := make(chan Item[int])

		go func() {
			defer close(in)
			in <- Item[int]{Val: 1}
			in <- Item[int]{Val: 2}
			cancel()
			in <- Item[int]{Val: 3}
			in <- Item[int]{Val: 4}
			in <- Item[int]{Val: 5}
		}()

		b := Batch[int](2, time.Second)

		got := Collect(b(ctx, in))

		assert.LessOrEqual(t, len(got), 3)
		assert.Equal(t, context.Canceled, got[len(got)-1].Err)
	})

	t.Run("with buffer size", func(t *testing.T) {
		ctx := context.Background()

		in := Of(1, 2, 3)

		b := Batch[int](2, time.Second, WithBufferSize(3))

		out := Pipe(in, b)(ctx, nil)

		got := Collect(out)

		want := []Item[[]int]{
			{Val: []int{1, 2}},
			{Val: []int{3}},
		}

		assert.Equal(t, 3, cap(out))
		assert.Equal(t, want, got)
	})

	t.Run("with pool size", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("with stop on error", func(t *testing.T) {
		ctx := context.Background()

		in := make(chan Item[int])

		go func() {
			defer close(in)
			in <- Item[int]{Val: 1}
			in <- Item[int]{Err: fmt.Errorf("error")}
			in <- Item[int]{Val: 2}
		}()

		b := Batch[int](2, time.Second, WithStopOnError(true))

		got := Collect(b(ctx, in))

		want := []Item[[]int]{
			{Err: fmt.Errorf("error")},
			{Val: []int{1}},
		}

		assert.Equal(t, want, got)
	})
}
