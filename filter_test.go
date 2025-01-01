package rivo_test

import (
	"context"
	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter(t *testing.T) {
	t.Run("filter all items", func(t *testing.T) {
		filterFn := func(ctx context.Context, i Item[int]) (bool, error) {
			return i.Val%2 == 0, nil
		}

		g := Of(1, 2, 3, 4, 5)

		f := Filter(filterFn)

		got := Pipe(g, f).Collect()

		want := []Item[int]{
			{Val: 2},
			{Val: 4},
		}

		assert.Equal(t, want, got)
	})

	t.Run("filter with error", func(t *testing.T) {
		filterFn := func(ctx context.Context, i Item[int]) (bool, error) {
			if i.Val == 3 {
				return false, assert.AnError
			}
			return i.Val%2 == 0, nil
		}

		g := Of(1, 2, 3, 4, 5)

		f := Filter(filterFn)

		got := Pipe(g, f).Collect()

		want := []Item[int]{
			{Val: 2},
			{Err: assert.AnError},
			{Val: 4},
		}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		filterFn := func(ctx context.Context, i Item[int]) (bool, error) {
			if i.Err != nil {
				return false, i.Err
			}
			return i.Val%2 == 0, nil
		}

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

		f := Filter(filterFn)

		got := Collect(f(ctx, in))

		assert.LessOrEqual(t, len(got), 4)
		assert.Equal(t, context.Canceled, got[len(got)-1].Err)
	})

	t.Run("with buffer size", func(t *testing.T) {
		filterFn := func(ctx context.Context, i Item[int]) (bool, error) {
			return i.Val%2 == 0, nil
		}

		in := make(chan Item[int])

		go func() {
			defer close(in)
			in <- Item[int]{Val: 1}
			in <- Item[int]{Val: 2}
			in <- Item[int]{Val: 3}
		}()

		f := Filter(filterFn, WithBufferSize(3))

		out := f(context.Background(), in)

		got := Collect(out)

		want := []Item[int]{
			{Val: 2},
		}

		assert.Equal(t, 3, cap(out))
		assert.Equal(t, want, got)
	})

	t.Run("with pool size", func(t *testing.T) {
		filterFn := func(ctx context.Context, i Item[int]) (bool, error) {
			return i.Val%2 == 0, nil
		}

		in := Of(1, 2, 3, 4, 5)

		f := Filter(filterFn, WithPoolSize(3))

		got := Pipe(in, f).Collect()

		want := []Item[int]{
			{Val: 2},
			{Val: 4},
		}

		assert.ElementsMatch(t, want, got)
	})

	t.Run("with stop on error", func(t *testing.T) {
		filterFn := func(ctx context.Context, i Item[int]) (bool, error) {
			if i.Val == 3 {
				return false, assert.AnError
			}
			return i.Val%2 == 0, nil
		}

		in := Of(1, 2, 3, 4, 5)

		f := Filter(filterFn, WithStopOnError(true))

		got := Pipe(in, f).Collect()

		want := []Item[int]{
			{Val: 2},
			{Err: assert.AnError},
		}

		assert.Equal(t, want, got)
	})
}
