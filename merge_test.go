package rivo_test

import (
	"context"
	"testing"
	"time"

	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	t.Run("merge two channels", func(t *testing.T) {
		ctx := context.Background()
		ch1 := make(chan int, 2)
		ch2 := make(chan int, 2)

		ch1 <- 1
		ch1 <- 2
		close(ch1)

		ch2 <- 3
		ch2 <- 4
		close(ch2)

		out := Merge(ctx, ch1, ch2)
		got := Collect(out)

		assert.ElementsMatch(t, []int{1, 2, 3, 4}, got)
	})

	t.Run("merge with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		ch1 := make(chan int)
		ch2 := make(chan int)

		go func() {
			defer close(ch1)
			for i := 0; i < 10; i++ {
				ch1 <- i
			}
		}()
		go func() {
			defer close(ch2)
			for i := 10; i < 20; i++ {
				ch2 <- i
			}
		}()

		// Cancel context quickly
		cancel()
		out := Merge(ctx, ch1, ch2)
		got := Collect(out)

		assert.Less(t, len(got), 5, "expected few or no items due to context cancellation")
	})

	t.Run("merge with one closed channel", func(t *testing.T) {
		ctx := context.Background()
		ch1 := make(chan int, 2)
		ch2 := make(chan int, 2)

		ch1 <- 1
		close(ch1)
		// ch2 is closed and empty
		close(ch2)

		out := Merge(ctx, ch1, ch2)
		got := Collect(out)

		assert.ElementsMatch(t, []int{1}, got)
	})

	t.Run("merge with slow channel", func(t *testing.T) {
		ctx := context.Background()
		ch1 := make(chan int, 1)
		ch2 := make(chan int, 1)

		go func() {
			defer close(ch1)
			ch1 <- 1
			time.Sleep(10 * time.Millisecond)
			ch1 <- 2
		}()
		go func() {
			defer close(ch2)
			ch2 <- 3
		}()

		out := Merge(ctx, ch1, ch2)
		got := Collect(out)

		assert.ElementsMatch(t, []int{1, 2, 3}, got)
	})
}
