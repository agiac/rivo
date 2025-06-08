package rivo_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	. "github.com/agiac/rivo"
	"github.com/agiac/rivo/core"
	"github.com/stretchr/testify/assert"
)

func ExampleFromFunc() {
	ctx := context.Background()

	count := atomic.Int32{}

	genFn := func(ctx context.Context) (int32, error) {
		value := count.Add(1)

		if value > 5 {
			return 0, ErrEOS
		}

		return value, nil
	}

	in := FromFunc(genFn)

	s := in(ctx, nil)

	for item := range s {
		fmt.Println(item.Val)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func TestFromFunc(t *testing.T) {
	t.Run("generate items from function", func(t *testing.T) {
		ctx := context.Background()
		count := 0
		genFn := func(ctx context.Context) (int, error) {
			count++
			if count > 5 {
				return 0, ErrEOS
			}
			return count, nil
		}

		f := FromFunc(genFn)

		got := core.Collect(f(ctx, nil))

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
		}

		assert.Equal(t, want, got)
	})

	t.Run("generate items with error", func(t *testing.T) {
		ctx := context.Background()
		count := 0
		genFn := func(ctx context.Context) (int, error) {
			count++
			if count == 3 {
				return 0, assert.AnError
			}
			if count > 5 {
				return 0, ErrEOS
			}
			return count, nil
		}

		f := FromFunc(genFn)

		got := core.Collect(f(ctx, nil))

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Err: assert.AnError},
			{Val: 4},
			{Val: 5},
		}

		assert.Equal(t, want, got)
	})
}
