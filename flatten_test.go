package rivo_test

import (
	"context"
	"fmt"
	"testing"

	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

func ExampleFlatten() {
	ctx := context.Background()

	in := Of([]int{1, 2}, []int{3, 4}, []int{5})

	f := Flatten[int]()

	p := Pipe(in, f)

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

		got := Collect(Pipe(in, f)(ctx, nil))

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

		got := Collect(f(ctx, in))

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
		defer cancel()

		in := make(chan Item[[]int])

		go func() {
			defer close(in)
			in <- Item[[]int]{Val: []int{1, 2}}
			cancel()
			in <- Item[[]int]{Val: []int{3, 4}}
			in <- Item[[]int]{Val: []int{5, 6}}
		}()

		f := Flatten[int]()

		got := Collect(f(ctx, in))

		assert.LessOrEqual(t, len(got), 4)
		assert.Equal(t, context.Canceled, got[len(got)-1].Err)
	})
}
