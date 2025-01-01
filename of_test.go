package rivo_test

import (
	"context"
	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOf(t *testing.T) {
	t.Run("create stream from items", func(t *testing.T) {
		got := Of(1, 2, 3, 4, 5).Collect()

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
		}

		assert.Equal(t, want, got)
	})

	t.Run("context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		in := Of(1, 2, 3, 4, 5)(ctx, nil)

		first := <-in
		second := <-in
		third := <-in
		cancel()
		fourth := <-in
		fifth, ok := <-in

		assert.Equal(t, Item[int]{Val: 1}, first)
		assert.Equal(t, Item[int]{Val: 2}, second)
		assert.Equal(t, Item[int]{Val: 3}, third)
		assert.Equal(t, Item[int]{Err: context.Canceled}, fourth)
		assert.Equal(t, Item[int]{}, fifth)
		assert.False(t, ok)
	})
}
