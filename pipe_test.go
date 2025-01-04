package rivo_test

import (
	"context"
	"fmt"
	"testing"

	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

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

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		})

		p := Pipe(a, b)

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
			return i.Val + 1, nil
		})

		p := Pipe(a, b)

		got := Collect(p(ctx, nil))

		assert.LessOrEqual(t, len(got), 4)
	})
}

func TestPipe2(t *testing.T) {
	t.Run("pipe all values", func(t *testing.T) {
		ctx := context.Background()

		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		})

		p := Pipe2(a, b)

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

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		})

		c := Filter(func(ctx context.Context, i Item[int]) (bool, error) {
			return i.Val%2 == 0, nil
		})

		p := Pipe3(a, b, c)

		got := Collect(p(ctx, nil))

		want := []Item[int]{
			{Val: 2},
			{Val: 4},
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
			return i.Val + 1, nil
		})

		c := Filter(func(ctx context.Context, i Item[int]) (bool, error) {
			if i.Val == 3 {
				cancel()
			}
			return i.Val%2 == 0, nil
		})

		p := Pipe3(a, b, c)

		got := Collect(p(ctx, nil))

		assert.LessOrEqual(t, len(got), 3)
	})
}

func TestPipe4(t *testing.T) {
	t.Run("pipe all values", func(t *testing.T) {
		ctx := context.Background()

		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		})

		c := Filter(func(ctx context.Context, i Item[int]) (bool, error) {
			return i.Val%2 == 0, nil
		})

		d := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val * 2, nil
		})

		p := Pipe4(a, b, c, d)

		got := Collect(p(ctx, nil))

		want := []Item[int]{
			{Val: 4},
			{Val: 8},
			{Val: 12},
		}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx := context.Background()

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		})

		c := Filter(func(ctx context.Context, i Item[int]) (bool, error) {
			if i.Val == 4 {
				cancel()
			}
			return i.Val%2 == 0, nil
		})

		d := Map(func(ctx context.Context, i Item[int]) (int, error) {
			if i.Err != nil {
				return 0, i.Err
			}
			return i.Val * 2, nil
		})

		p := Pipe4(a, b, c, d)

		got := Collect(p(ctx, nil))

		assert.LessOrEqual(t, len(got), 3)
	})
}

func TestPipe5(t *testing.T) {
	t.Run("pipe all values", func(t *testing.T) {
		ctx := context.Background()

		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		})

		c := Filter(func(ctx context.Context, i Item[int]) (bool, error) {
			return i.Val%2 == 0, nil
		})

		d := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val * 2, nil
		})

		e := Filter(func(ctx context.Context, i Item[int]) (bool, error) {
			return i.Val%3 == 0, nil
		})

		p := Pipe5(a, b, c, d, e)

		got := Collect(p(ctx, nil))

		want := []Item[int]{
			{Val: 12},
		}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx := context.Background()

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		})

		c := Filter(func(ctx context.Context, i Item[int]) (bool, error) {
			if i.Val == 4 {
				cancel()
			}
			return i.Val%2 == 0, nil
		})

		d := Map(func(ctx context.Context, i Item[int]) (int, error) {
			if i.Err != nil {
				return 0, i.Err
			}
			return i.Val * 2, nil
		})

		e := Filter(func(ctx context.Context, i Item[int]) (bool, error) {
			return i.Val%3 == 0, nil
		})

		p := Pipe5(a, b, c, d, e)

		got := Collect(p(ctx, nil))

		assert.LessOrEqual(t, len(got), 2)
	})
}
