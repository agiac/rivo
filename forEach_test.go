package rivo_test

import (
	"context"
	"fmt"
	"testing"

	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

func ExampleForEach() {
	ctx := context.Background()

	g := Of(1, 2, 3, 4, 5)

	f := ForEach(func(ctx context.Context, i Item[int]) error {
		// Do some side effect
		// ...
		// Simulate an error
		if i.Val == 3 {
			return fmt.Errorf("an error")
		}

		return nil
	})

	s := Pipe(g, f)(ctx, nil)

	for item := range s {
		fmt.Printf("item: %v; error: %v\n", item.Val, item.Err)
	}

	// Output:
	// item: {}; error: an error
}

func TestForEach(t *testing.T) {
	t.Run("for each item", func(t *testing.T) {
		ctx := context.Background()

		g := Of(1, 2, 3, 4, 5)

		sideEffect := make([]int, 0)
		f := ForEach(func(ctx context.Context, i Item[int]) error {
			sideEffect = append(sideEffect, i.Val)
			return nil
		})

		errs := Collect(Pipe(g, f)(ctx, nil))

		assert.Equal(t, []int{1, 2, 3, 4, 5}, sideEffect)
		assert.Equal(t, 0, len(errs))
	})

	t.Run("forward errors", func(t *testing.T) {
		ctx := context.Background()

		g := Of(1, 2, 3, 4, 5)

		sideEffect := make([]int, 0)
		f := ForEach(func(ctx context.Context, i Item[int]) error {
			sideEffect = append(sideEffect, i.Val)

			if i.Val == 3 {
				return fmt.Errorf("an error")
			}

			return nil
		})

		errs := Collect(Pipe(g, f)(ctx, nil))

		assert.Equal(t, []int{1, 2, 3, 4, 5}, sideEffect)
		assert.Equal(t, 1, len(errs))
		assert.Error(t, errs[0].Err)
	})
}
