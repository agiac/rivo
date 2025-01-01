package rivo_test

import (
	"context"
	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("map all items", func(t *testing.T) {
		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		}

		g := Of(1, 2, 3, 4, 5)

		m := Map(mapFn)

		got := Pipe(g, m).Collect()

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
			{Val: 6},
		}

		assert.Equal(t, want, got)
	})

	t.Run("map all items with error", func(t *testing.T) {
		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			if i.Val == 3 {
				return 0, assert.AnError
			}
			return i.Val + 1, nil
		}

		g := Of(1, 2, 3, 4, 5)

		m := Map(mapFn)

		got := Pipe(g, m).Collect()

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Err: assert.AnError},
			{Val: 5},
			{Val: 6},
		}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			if i.Err != nil {
				return 0, i.Err
			}

			return i.Val + 1, nil
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

		m := Map(mapFn)

		got := Collect(m(ctx, in))

		assert.LessOrEqual(t, len(got), 4)
		assert.Equal(t, context.Canceled, got[len(got)-1].Err)
	})

	t.Run("with buffer size", func(t *testing.T) {
		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		}

		in := make(chan Item[int])

		go func() {
			defer close(in)
			in <- Item[int]{Val: 1}
			in <- Item[int]{Val: 2}
			in <- Item[int]{Val: 3}
		}()

		m := Map(mapFn, WithBufferSize(3))

		out := m(context.Background(), in)

		got := Collect(out)

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Val: 4},
		}

		assert.Equal(t, 3, cap(out))
		assert.Equal(t, want, got)
	})

	t.Run("with pool size", func(t *testing.T) {
		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		}

		in := Of(1, 2, 3, 4, 5)

		m := Map(mapFn, WithPoolSize(3))

		got := Pipe(in, m).Collect()

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
			{Val: 6},
		}

		assert.ElementsMatch(t, want, got)
	})

	t.Run("with stop on error", func(t *testing.T) {
		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			if i.Val == 3 {
				return 0, assert.AnError
			}
			return i.Val + 1, nil
		}

		in := Of(1, 2, 3, 4, 5)

		m := Map(mapFn, WithStopOnError(true))

		got := Pipe(in, m).Collect()

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Err: assert.AnError},
		}

		assert.Equal(t, want, got)
	})
}
