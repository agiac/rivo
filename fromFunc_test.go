package rivo_test

import (
	"context"
	"fmt"
	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
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

		got := Collect(f(ctx, nil))

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

		got := Collect(f(ctx, nil))

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Err: assert.AnError},
			{Val: 4},
			{Val: 5},
		}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		count := 0
		genFn := func(ctx context.Context) (int, error) {
			count++
			if count > 2 {
				cancel()
			}
			return count, nil
		}

		f := FromFunc(genFn)

		got := Collect(f(ctx, nil))

		assert.LessOrEqual(t, len(got), 3)
		assert.Equal(t, context.Canceled, got[len(got)-1].Err)
	})

	t.Run("with buffer size", func(t *testing.T) {
		ctx := context.Background()
		count := 0
		genFn := func(ctx context.Context) (int, error) {
			count++
			if count > 3 {
				return 0, ErrEOS
			}
			return count, nil
		}

		f := FromFunc(genFn, WithBufferSize(3))

		in := make(chan Item[struct{}])

		go func() {
			defer close(in)
			in <- Item[struct{}]{}
			in <- Item[struct{}]{}
			in <- Item[struct{}]{}
		}()

		out := f(ctx, in)

		got := Collect(out)

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Val: 3},
		}

		assert.Equal(t, 3, cap(out))
		assert.Equal(t, want, got)
	})

	t.Run("with pool size", func(t *testing.T) {
		ctx := context.Background()
		count := atomic.Int32{}
		genFn := func(ctx context.Context) (int, error) {
			v := int(count.Add(1))
			if v > 5 {
				return 0, ErrEOS
			}
			return v, nil
		}

		f := FromFunc(genFn, WithPoolSize(3))

		got := Collect(f(ctx, nil))

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
		}

		assert.ElementsMatch(t, want, got)
	})

	t.Run("with stop on error", func(t *testing.T) {
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

		f := FromFunc(genFn, WithStopOnError(true))

		got := Collect(f(ctx, nil))

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Err: assert.AnError},
		}

		assert.Equal(t, want, got)
	})
}
