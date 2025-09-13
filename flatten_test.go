package rivo_test

import (
	"context"
	"fmt"
	. "github.com/agiac/rivo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleFlatten() {
	ctx := context.Background()

	in := Of([]int{1, 2}, []int{3, 4}, []int{5})

	f := Flatten[int]()

	p := Pipe(in, f)

	for item := range p(ctx, nil, nil) {
		fmt.Printf("%v\n", item)
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

		got := Collect(Pipe(in, f)(ctx, nil, nil))

		want := []int{1, 2, 3, 4, 5}

		assert.Equal(t, want, got)
	})

	t.Run("context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		g := Of([]int{1, 2}, []int{3, 4}, []int{5})
		f := Flatten[int]()

		got := Collect(Pipe(g, f)(ctx, nil, nil))

		assert.Lessf(t, len(got), 3, "expected less than 3 items due to context cancellation")
	})
}
