package rivo_test

import (
	"context"
	"fmt"
	"testing"

	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

var addOne = Map(func(ctx context.Context, i Item[int]) (int, error) {
	return i.Val + 1, nil
})

func ExamplePipe() {
	ctx := context.Background()

	a := Of(1, 2, 3, 4, 5)

	b := Map(func(ctx context.Context, i Item[int]) (int, error) {
		return i.Val + 1, nil
	})

	p := Pipe(a, b)

	s := p(ctx, nil)

	for item := range s {
		fmt.Println(item.Val)
	}

	// Output:
	// 2
	// 3
	// 4
	// 5
	// 6
}

func TestPipe(t *testing.T) {
	t.Run("pipe all values", func(t *testing.T) {
		ctx := context.Background()

		a := Of(1, 2, 3, 4, 5)

		p := Pipe(a, addOne)

		got := Collect(p(ctx, nil))

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
			{Val: 6},
		}

		assert.Equal(t, want, got)
	})

	t.Run("context cancelled", func(t *testing.T) {
		ctx := context.Background()

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			if i.Val == 2 {
				cancel()
			}

			if i.Err != nil {
				return 0, i.Err
			}

			return i.Val + 1, nil
		})

		p := Pipe(a, b)

		got := Collect(p(ctx, nil))

		assert.LessOrEqual(t, len(got), 5)
		assert.Equal(t, ctx.Err(), got[len(got)-1].Err)
	})
}

func TestPipe2(t *testing.T) {
	t.Run("pipe all values", func(t *testing.T) {
		ctx := context.Background()

		a := Of(1, 2, 3, 4, 5)

		p := Pipe2(a, addOne)

		got := Collect(p(ctx, nil))

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
			{Val: 6},
		}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx := context.Background()

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			if i.Val == 2 {
				cancel()
			}
			return i.Val + 1, nil
		})

		p := Pipe2(a, b)

		got := Collect(p(ctx, nil))

		assert.LessOrEqual(t, len(got), 4)
	})
}

func TestPipe3(t *testing.T) {
	t.Run("pipe all values", func(t *testing.T) {
		ctx := context.Background()

		a := Of(1, 2, 3, 4, 5)

		p := Pipe3(a, addOne, addOne)

		got := Collect(p(ctx, nil))

		want := []Item[int]{
			{Val: 3},
			{Val: 4},
			{Val: 5},
			{Val: 6},
			{Val: 7},
		}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx := context.Background()

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		a := Of(1, 2, 3, 4, 5)

		f := Filter(func(ctx context.Context, i Item[int]) (bool, error) {
			if i.Val == 3 {
				cancel()
			}
			return i.Val%2 == 0, nil
		})

		p := Pipe3(a, addOne, f)

		got := Collect(p(ctx, nil))

		assert.LessOrEqual(t, len(got), 3)
	})
}

func TestPipe4(t *testing.T) {
	t.Run("pipe all values", func(t *testing.T) {
		ctx := context.Background()

		a := Of(1, 2, 3, 4, 5)

		p := Pipe4(a, addOne, addOne, addOne)

		got := Collect(p(ctx, nil))

		want := []Item[int]{
			{Val: 4},
			{Val: 5},
			{Val: 6},
			{Val: 7},
			{Val: 8},
		}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx := context.Background()

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		a := Of(1, 2, 3, 4, 5)

		f := Filter(func(ctx context.Context, i Item[int]) (bool, error) {
			if i.Val == 4 {
				cancel()
			}
			return i.Val%2 == 0, nil
		})

		p := Pipe4(a, addOne, f, addOne)

		got := Collect(p(ctx, nil))

		assert.LessOrEqual(t, len(got), 3)
	})
}

func TestPipe5(t *testing.T) {
	t.Run("pipe all values", func(t *testing.T) {
		ctx := context.Background()

		a := Of(1, 2, 3, 4, 5)

		p := Pipe5(a, addOne, addOne, addOne, addOne)

		got := Collect(p(ctx, nil))

		want := []Item[int]{
			{Val: 5},
			{Val: 6},
			{Val: 7},
			{Val: 8},
			{Val: 9},
		}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx := context.Background()

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		a := Of(1, 2, 3, 4, 5)

		f := Filter(func(ctx context.Context, i Item[int]) (bool, error) {
			if i.Val == 4 {
				cancel()
			}
			return i.Val%2 == 0, nil
		})

		p := Pipe5(a, addOne, f, addOne, addOne)

		got := Collect(p(ctx, nil))

		assert.LessOrEqual(t, len(got), 3)
	})
}
