package rivo_test

import (
	"context"
	. "github.com/agiac/rivo"
	"github.com/agiac/rivo/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrDone(t *testing.T) {
	t.Run("continue till end of input stream", func(t *testing.T) {
		ctx := context.Background()

		in := Of(1, 2, 3, 4, 5)(ctx, nil)

		got := core.Collect(OrDone(ctx, in))

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
		}

		assert.Equal(t, want, got)
	})

	//t.Run("continue till context is cancelled", func(t *testing.T) {
	//	ctx, cancel := context.WithCancel(context.Background())
	//
	//	in := make(chan Item[int])
	//
	//	go func() {
	//		in <- Item[int]{Val: 1}
	//		in <- Item[int]{Val: 2}
	//		cancel()
	//		in <- Item[int]{Val: 3}
	//		in <- Item[int]{Val: 4}
	//		in <- Item[int]{Val: 5}
	//	}()
	//
	//	got := core.Collect(OrDone(ctx, in))
	//
	//	assert.LessOrEqual(t, len(got), 3)
	//	assert.ErrorIs(t, got[len(got)-1].Err, context.Canceled)
	//})
}
