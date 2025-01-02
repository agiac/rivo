package rivo_test

import (
	"context"
	"fmt"
	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExampleMap() {
	ctx := context.Background()

	in := Of(1, 2, 3, 4, 5)

	double := Map(func(ctx context.Context, i Item[int]) (int, error) {
		// Always check for errors
		if i.Err != nil {
			return 0, i.Err // Propagate the error
		}

		return i.Val * 2, nil
	})

	p := Pipe(in, double)

	s := p(ctx, nil)

	for item := range s {
		fmt.Println(item.Val)
	}

	// Output:
	// 2
	// 4
	// 6
	// 8
	// 10
}

func TestMap(t *testing.T) {
	t.Run("map all items", func(t *testing.T) {
		ctx := context.Background()

		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		}

		g := Of(1, 2, 3, 4, 5)

		m := Map(mapFn)

		got := Collect(Pipe(g, m)(ctx, nil))

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
		ctx := context.Background()

		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			if i.Val == 3 {
				return 0, assert.AnError
			}
			return i.Val + 1, nil
		}

		g := Of(1, 2, 3, 4, 5)

		m := Map(mapFn)

		got := Collect(Pipe(g, m)(ctx, nil))

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
		ctx := context.Background()

		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			if i.Err != nil {
				return 0, i.Err
			}

			return i.Val + 1, nil
		}

		ctx, cancel := context.WithCancel(ctx)
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
		ctx := context.Background()

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

		out := m(ctx, in)

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
		ctx := context.Background()

		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		}

		in := Of(1, 2, 3, 4, 5)

		m := Map(mapFn, WithPoolSize(3))

		got := Collect(Pipe(in, m)(ctx, nil))

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
		ctx := context.Background()

		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			if i.Val == 3 {
				return 0, assert.AnError
			}
			return i.Val + 1, nil
		}

		in := Of(1, 2, 3, 4, 5)

		m := Map(mapFn, WithStopOnError(true))

		got := Collect(Pipe(in, m)(ctx, nil))

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Err: assert.AnError},
		}

		assert.Equal(t, want, got)
	})

	t.Run("with on before close", func(t *testing.T) {
		ctx := context.Background()

		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		}

		in := Of(1, 2, 3, 4, 5)

		beforeCloseCalled := false

		m := Map(mapFn, WithOnBeforeClose(func(ctx context.Context) error {
			beforeCloseCalled = true
			return nil
		}))

		_ = Collect(Pipe(in, m)(ctx, nil))

		assert.True(t, beforeCloseCalled)
	})

	t.Run("with on before close error", func(t *testing.T) {
		ctx := context.Background()

		mapFn := func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		}

		in := Of(1, 2, 3, 4, 5)

		m := Map(mapFn, WithOnBeforeClose(func(ctx context.Context) error {
			return assert.AnError
		}))

		got := Collect(Pipe(in, m)(ctx, nil))

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
			{Val: 6},
			{Err: assert.AnError},
		}

		assert.Equal(t, want, got)
	})
}
