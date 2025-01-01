package rivo_test

import (
	"context"
	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollect(t *testing.T) {
	t.Run("collect till end of input stream", func(t *testing.T) {
		ctx := context.Background()

		in := Of(1, 2, 3, 4, 5)(ctx, nil)

		got := Collect(in)

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
		}

		assert.Equal(t, want, got)
	})
}

func TestCollectWithContext(t *testing.T) {
	t.Run("collect till end of input stream", func(t *testing.T) {
		ctx := context.Background()

		in := Of(1, 2, 3, 4, 5)(ctx, nil)

		got := CollectWithContext(ctx, in)

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
		}

		assert.Equal(t, want, got)
	})

	t.Run("collect till context is cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		in := make(chan Item[int], 5)

		go func() {
			defer close(in)
			in <- Item[int]{Val: 1}
			in <- Item[int]{Val: 2}
			cancel()
			in <- Item[int]{Val: 3}
			in <- Item[int]{Val: 4}
			in <- Item[int]{Val: 5}
		}()

		got := CollectWithContext(ctx, in)

		assert.LessOrEqual(t, len(got), 4)
		assert.Equal(t, context.Canceled, got[len(got)-1].Err)
	})
}
