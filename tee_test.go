package rivo_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

func ExampleTee2() {
	ctx := context.Background()

	g := rivo.Of("hello", "hello", "hello")

	out1, out2 := rivo.Tee2(g)(ctx, nil)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := range out1(ctx, nil) {
			fmt.Println(i.Val)
		}
	}()

	go func() {
		defer wg.Done()
		for i := range out2(ctx, nil) {
			fmt.Println(i.Val)
		}
	}()

	wg.Wait()

	// Output:
	// hello
	// hello
	// hello
	// hello
	// hello
	// hello
}

func TestTee2(t *testing.T) {
	t.Run("tee stream", func(t *testing.T) {
		ctx := context.Background()

		g := rivo.Of("hello", "hello", "hello")

		out1, out2 := rivo.Tee2(g)(ctx, nil)

		var got1, got2 []rivo.Item[string]
		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			got1 = rivo.Collect(out1(ctx, nil))
		}()

		go func() {
			defer wg.Done()
			got2 = rivo.Collect(out2(ctx, nil))
		}()

		wg.Wait()

		want := []rivo.Item[string]{
			{Val: "hello"},
			{Val: "hello"},
			{Val: "hello"},
		}

		assert.Equal(t, want, got1)
		assert.Equal(t, want, got2)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		in := make(chan rivo.Item[string])
		go func() {
			defer close(in)
			in <- rivo.Item[string]{Val: "hello"}
			cancel()
			in <- rivo.Item[string]{Val: "hello"}
			in <- rivo.Item[string]{Val: "hello"}
			in <- rivo.Item[string]{Val: "hello"}
		}()

		g := func(ctx context.Context, s rivo.Stream[string]) rivo.Stream[string] {
			return in
		}

		out1, out2 := rivo.Tee2(g)(ctx, nil)

		var got1, got2 []rivo.Item[string]
		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			got1 = rivo.Collect(out1(ctx, nil))
		}()

		go func() {
			defer wg.Done()
			got2 = rivo.Collect(out2(ctx, nil))
		}()

		wg.Wait()

		assert.LessOrEqual(t, len(got1), 3)
		assert.Equal(t, context.Canceled, got1[len(got1)-1].Err)
		assert.LessOrEqual(t, len(got2), 3)
		assert.Equal(t, context.Canceled, got2[len(got2)-1].Err)
	})
}

func TestTee2N(t *testing.T) {
	t.Run("tee stream", func(t *testing.T) {
		ctx := context.Background()

		g := rivo.Of("hello", "hello", "hello")

		const n = 5

		out := rivo.Tee2N(g, n)(ctx, nil)

		got := make([][]rivo.Item[string], n)
		wg := sync.WaitGroup{}
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func(i int) {
				defer wg.Done()
				got[i] = rivo.Collect(out[i](ctx, nil))
			}(i)
		}

		wg.Wait()

		want := []rivo.Item[string]{
			{Val: "hello"},
			{Val: "hello"},
			{Val: "hello"},
		}

		for i := 0; i < n; i++ {
			assert.Equal(t, want, got[i])
		}
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		in := make(chan rivo.Item[string])
		go func() {
			defer close(in)
			in <- rivo.Item[string]{Val: "hello"}
			cancel()
			in <- rivo.Item[string]{Val: "hello"}
			in <- rivo.Item[string]{Val: "hello"}
			in <- rivo.Item[string]{Val: "hello"}
		}()

		const n = 5

		out := rivo.Tee2N(func(ctx context.Context, s rivo.Stream[string]) rivo.Stream[string] {
			return in
		}, n)(ctx, nil)

		got := make([][]rivo.Item[string], n)
		wg := sync.WaitGroup{}
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func(i int) {
				defer wg.Done()
				got[i] = rivo.Collect(out[i](ctx, nil))
			}(i)
		}

		wg.Wait()

		for i := 0; i < n; i++ {
			assert.LessOrEqual(t, len(got[i]), 3)
			assert.Equal(t, context.Canceled, got[i][len(got[i])-1].Err)
		}
	})
}
