package rivo_test

import (
	"context"
	"fmt"
	. "github.com/agiac/rivo"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleFromFunc() {
	ctx := context.Background()

	count := atomic.Int32{}

	genFn := func(ctx context.Context) (int32, bool) {
		value := count.Add(1)

		if value > 5 {
			return 0, false
		}

		return value, true
	}

	in := FromFunc(genFn)

	s := in(ctx, nil)

	for item := range s {
		fmt.Println(item)
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
		genFn := func(ctx context.Context) (int, bool) {
			count++
			if count > 5 {
				return 0, false
			}
			return count, true
		}

		f := FromFunc(genFn)

		got := Collect(f(ctx, nil))

		want := []int{1, 2, 3, 4, 5}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		count := 0
		genFn := func(ctx context.Context) (int, bool) {
			count++
			if count > 5 {
				return 0, false
			}
			return count, true
		}

		f := FromFunc(genFn)

		got := Collect(f(ctx, nil))

		assert.Less(t, len(got), 5, "should not generate more than 5 items when context is cancelled")
	})

	t.Run("with buffer size", func(t *testing.T) {
		ctx := context.Background()
		count := 0
		genFn := func(ctx context.Context) (int, bool) {
			count++
			if count > 5 {
				return 0, false
			}
			return count, true
		}

		f := FromFunc(genFn, FromFuncBufferSize(3))

		out := f(ctx, nil)

		got := Collect(out)

		want := []int{1, 2, 3, 4, 5}

		assert.Equal(t, 3, cap(out))
		assert.Equal(t, want, got)
	})

	t.Run("with pool size", func(t *testing.T) {
		ctx := context.Background()
		count := atomic.Int32{}
		genFn := func(ctx context.Context) (int, bool) {
			v := int(count.Add(1))
			if v > 5 {
				return 0, false
			}
			return v, true
		}

		f := FromFunc(genFn, FromFuncPoolSize(3))

		got := Collect(f(ctx, nil))

		want := []int{1, 2, 3, 4, 5}

		assert.ElementsMatch(t, want, got)
	})

	t.Run("with on before close", func(t *testing.T) {
		ctx := context.Background()
		var n int
		genFn := func(ctx context.Context) (int, bool) {
			n++
			if n > 5 {
				return 0, false
			}
			return n, true
		}

		beforeCloseCalled := atomic.Bool{}

		f := FromFunc(genFn, FromFuncOnBeforeClose(func(ctx context.Context) {
			beforeCloseCalled.Store(true)
		}))

		got := Collect(f(ctx, nil))

		want := []int{1, 2, 3, 4, 5}

		assert.Equal(t, want, got)
		assert.True(t, beforeCloseCalled.Load())
	})
}
