package rivo_test

import (
	"context"
	"fmt"
	"testing"

	. "github.com/agiac/rivo"

	"github.com/stretchr/testify/assert"
)

func ExampleForEachOutput() {
	ctx := context.Background()

	in := Of(1, 2, 3, 4, 5)

	f := func(ctx context.Context, n int, out chan<- int, errs chan<- error) {
		out <- n * 2
	}

	p := Pipe(in, ForEachOutput(f))

	s := p(ctx, nil, nil)

	for n := range s {
		fmt.Println(n)
	}

	// Output:
	// 2
	// 4
	// 6
	// 8
	// 10
}

func TestForEachOutput(t *testing.T) {
	t.Run("for each output all items", func(t *testing.T) {
		ctx := context.Background()

		f := func(ctx context.Context, n int, out chan<- int, errs chan<- error) {
			out <- n + 1
		}

		g := Of(1, 2, 3, 4, 5)

		fo := ForEachOutput(f)

		got := Collect(Pipe(g, fo)(ctx, nil, nil))

		want := []int{2, 3, 4, 5, 6}

		assert.ElementsMatch(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx := context.Background()

		f := func(ctx context.Context, n int, out chan<- int, errs chan<- error) {
			out <- n + 1
		}

		ctx, cancel := context.WithCancel(ctx)
		cancel()

		g := Of(1, 2, 3, 4, 6)
		fo := ForEachOutput(f)

		got := Collect(Pipe(g, fo)(ctx, nil, nil))

		assert.Lessf(t, len(got), 3, "expected less than 3 items due to context cancellation")
	})

	t.Run("with buffer size", func(t *testing.T) {
		ctx := context.Background()

		f := func(ctx context.Context, n int, out chan<- int, errs chan<- error) {
			out <- n + 1
		}

		in := make(chan int)

		go func() {
			defer close(in)
			in <- 1
			in <- 2
			in <- 3
		}()

		fo := ForEachOutput(f, ForEachOutputBufferSize(3))

		out := fo(ctx, in, nil)

		got := Collect(out)

		want := []int{2, 3, 4}

		assert.Equal(t, 3, cap(out))
		assert.ElementsMatch(t, want, got)
	})

	t.Run("with pool size", func(t *testing.T) {
		ctx := context.Background()

		f := func(ctx context.Context, n int, out chan<- int, errs chan<- error) {
			out <- n + 1
		}

		in := Of(1, 2, 3, 4, 5)

		fo := ForEachOutput(f, ForEachOutputPoolSize(3))

		got := Collect(Pipe(in, fo)(ctx, nil, nil))

		want := []int{2, 3, 4, 5, 6}

		assert.ElementsMatch(t, want, got)
	})
}
