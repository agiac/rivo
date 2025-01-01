package rivo_test

import (
	"context"
	"fmt"
	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExamplePipe() {
	a := Of(1, 2, 3, 4, 5)

	b := Map(func(ctx context.Context, i Item[int]) (int, error) {
		return i.Val + 1, nil
	})

	p := Pipe(a, b)

	for _, item := range p.Collect() {
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
	t.Run("pipe and collect", func(t *testing.T) {
		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		})

		got := Pipe(a, b).Collect()

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
			{Val: 6},
		}

		assert.Equal(t, want, got)
	})

	t.Run("pipe and collect with context", func(t *testing.T) {
		ctx := context.Background()

		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		})

		got := Pipe(a, b).CollectWithContext(ctx)

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
			{Val: 6},
		}

		assert.Equal(t, want, got)
	})

	t.Run("pipe and collect with context and cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			if i.Val == 3 {
				cancel()
			}
			return i.Val + 1, nil
		})

		got := Pipe(a, b).CollectWithContext(ctx)

		assert.LessOrEqual(t, len(got), 3)
		assert.ErrorIs(t, got[len(got)-1].Err, context.Canceled)
	})
}

func TestPipe2(t *testing.T) {
	t.Run("pipe and collect", func(t *testing.T) {
		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		})

		got := Pipe2(a, b).Collect()

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
			{Val: 6},
		}

		assert.Equal(t, want, got)
	})

	t.Run("pipe and collect with context", func(t *testing.T) {
		ctx := context.Background()

		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		})

		got := Pipe2(a, b).CollectWithContext(ctx)

		want := []Item[int]{
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
			{Val: 6},
		}

		assert.Equal(t, want, got)
	})

	t.Run("pipe and collect with context and cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			if i.Val == 3 {
				cancel()
			}
			return i.Val + 1, nil
		})

		got := Pipe2(a, b).CollectWithContext(ctx)

		assert.LessOrEqual(t, len(got), 3)
		assert.ErrorIs(t, got[len(got)-1].Err, context.Canceled)
	})
}

func TestPipe3(t *testing.T) {
	t.Run("pipe and collect", func(t *testing.T) {
		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		})

		c := Filter(func(ctx context.Context, i Item[int]) (bool, error) {
			return i.Val%2 == 0, nil
		})

		got := Pipe3(a, b, c).Collect()

		want := []Item[int]{
			{Val: 2},
			{Val: 4},
			{Val: 6},
		}

		assert.Equal(t, want, got)
	})

	t.Run("pipe and collect with context", func(t *testing.T) {
		ctx := context.Background()

		a := Of(1, 2, 3, 4, 5)

		b := Map(func(ctx context.Context, i Item[int]) (int, error) {
			return i.Val + 1, nil
		})

		c := Filter(func(ctx context.Context, i Item[int]) (bool, error) {
			return i.Val%2 == 0, nil
		})

		got := Pipe3(a, b, c).CollectWithContext(ctx)

		want := []Item[int]{
			{Val: 2},
			{Val: 4},
			{Val: 6},
		}

		assert.Equal(t, want, got)
	})

	t.Run("pipe and collect with context and cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
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

		got := Pipe3(a, b, c).CollectWithContext(ctx)

		assert.LessOrEqual(t, len(got), 3)
		assert.ErrorIs(t, got[len(got)-1].Err, context.Canceled)
	})
}

func TestPipe4(t *testing.T) {
	t.Run("pipe and collect", func(t *testing.T) {
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

		got := Pipe4(a, b, c, d).Collect()

		want := []Item[int]{
			{Val: 4},
			{Val: 8},
			{Val: 12},
		}

		assert.Equal(t, want, got)
	})

	t.Run("pipe and collect with context", func(t *testing.T) {
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

		got := Pipe4(a, b, c, d).CollectWithContext(ctx)

		want := []Item[int]{
			{Val: 4},
			{Val: 8},
			{Val: 12},
		}

		assert.Equal(t, want, got)
	})

	t.Run("pipe and collect with context and cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
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
			return i.Val * 2, nil
		})

		got := Pipe4(a, b, c, d).CollectWithContext(ctx)

		assert.LessOrEqual(t, len(got), 3)
		assert.ErrorIs(t, got[len(got)-1].Err, context.Canceled)
	})
}

func TestPipe5(t *testing.T) {
	t.Run("pipe and collect", func(t *testing.T) {
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

		got := Pipe5(a, b, c, d, e).Collect()

		want := []Item[int]{
			{Val: 12},
		}

		assert.Equal(t, want, got)
	})

	t.Run("pipe and collect with context", func(t *testing.T) {
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

		got := Pipe5(a, b, c, d, e).CollectWithContext(ctx)

		want := []Item[int]{
			{Val: 12},
		}

		assert.Equal(t, want, got)
	})

	t.Run("pipe and collect with context and cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
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
			return i.Val * 2, nil
		})

		e := Filter(func(ctx context.Context, i Item[int]) (bool, error) {
			return i.Val%3 == 0, nil
		})

		got := Pipe5(a, b, c, d, e).CollectWithContext(ctx)

		assert.LessOrEqual(t, len(got), 2)
		assert.ErrorIs(t, got[len(got)-1].Err, context.Canceled)
	})
}
