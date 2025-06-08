package rivo_test

import (
	"context"
	"fmt"
	"github.com/agiac/rivo/core"
	"testing"

	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

func ExampleFlatten() {
	ctx := context.Background()

	in := Of([]int{1, 2}, []int{3, 4}, []int{5})

	f := Flatten[int]()

	p := core.Pipe(in, f)

	for item := range p(ctx, nil) {
		fmt.Printf("%v\n", item.Val)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func TestFlatten(t *testing.T) {
	t.Run("flatten slices", func(t *testing.T) {
		ctx := context.Background()

		in := Of([]int{1, 2}, []int{3, 4}, []int{5})

		f := Flatten[int]()

		got := core.Collect(core.Pipe(in, f)(ctx, nil))

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
		}

		assert.Equal(t, want, got)
	})

	t.Run("flatten with errors", func(t *testing.T) {
		ctx := context.Background()

		in := make(chan Item[[]int])

		go func() {
			defer close(in)
			in <- Item[[]int]{Val: []int{1, 2}}
			in <- Item[[]int]{Err: fmt.Errorf("error")}
			in <- Item[[]int]{Val: []int{3, 4}}
		}()

		f := Flatten[int]()

		got := core.Collect(f(ctx, in))

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Err: fmt.Errorf("error")},
			{Val: 3},
			{Val: 4},
		}

		assert.Equal(t, want, got)
	})

	t.Run("context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		g := Of([]int{1, 2}, []int{3, 4}, []int{5})
		f := Flatten[int]()

		got := core.Collect(core.Pipe(g, f)(ctx, nil))

		assert.Lessf(t, len(got), 3, "expected less than 3 items due to context cancellation")
	})
}
