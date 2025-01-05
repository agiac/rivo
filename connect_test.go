package rivo_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

func ExampleConnect() {
	ctx := context.Background()

	g := rivo.Of("Hello", "Hello", "Hello")

	capitalize := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (string, error) {
		return strings.ToUpper(i.Val), nil
	})

	lowercase := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (string, error) {
		return strings.ToLower(i.Val), nil
	})

	resA := make([]string, 0)
	a := rivo.Do(func(ctx context.Context, i rivo.Item[string]) {
		resA = append(resA, i.Val)
	})

	resB := make([]string, 0)
	b := rivo.Do(func(ctx context.Context, i rivo.Item[string]) {
		resB = append(resB, i.Val)
	})

	p1 := rivo.Pipe(capitalize, a)
	p2 := rivo.Pipe(lowercase, b)

	<-rivo.Connect(p1, p2)(ctx, g(ctx, nil))

	for _, s := range resA {
		fmt.Println(s)
	}

	for _, s := range resB {
		fmt.Println(s)
	}

	// Output:
	// HELLO
	// HELLO
	// HELLO
	// hello
	// hello
	// hello
}

func TestParallel(t *testing.T) {
	t.Run("run in parallel", func(t *testing.T) {
		ctx := context.Background()

		g := rivo.Of("Hello", "Hello", "Hello")

		resA := make([]string, 0)
		a := rivo.Do(func(ctx context.Context, i rivo.Item[string]) {
			resA = append(resA, i.Val)
		})

		resB := make([]string, 0)
		b := rivo.Do(func(ctx context.Context, i rivo.Item[string]) {
			resB = append(resB, i.Val)
		})

		<-rivo.Connect(a, b)(ctx, g(ctx, nil))

		want := []string{"Hello", "Hello", "Hello"}

		assert.Equal(t, want, resA)
		assert.Equal(t, want, resB)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		in := make(chan rivo.Item[string])
		go func() {
			defer close(in)
			in <- rivo.Item[string]{Val: "Hello"}
			in <- rivo.Item[string]{Val: "Hello"}
			cancel()
			in <- rivo.Item[string]{Val: "Hello"}
			in <- rivo.Item[string]{Val: "Hello"}
			in <- rivo.Item[string]{Val: "Hello"}

		}()

		mtxA := sync.Mutex{}
		resA := make([]rivo.Item[string], 0)
		a := rivo.Do(func(ctx context.Context, i rivo.Item[string]) {
			mtxA.Lock()
			defer mtxA.Unlock()
			resA = append(resA, i)
		})

		mtxB := sync.Mutex{}
		resB := make([]rivo.Item[string], 0)
		b := rivo.Do(func(ctx context.Context, i rivo.Item[string]) {
			mtxB.Lock()
			defer mtxB.Unlock()
			resB = append(resB, i)
		})

		<-rivo.Connect(a, b)(ctx, in)

		mtxA.Lock()
		mtxB.Lock()
		defer mtxA.Unlock()
		defer mtxB.Unlock()

		assert.LessOrEqual(t, len(resA), 3)
		assert.LessOrEqual(t, len(resB), 3)
		assert.Equal(t, context.Canceled, resA[len(resA)-1].Err)
		assert.Equal(t, context.Canceled, resB[len(resB)-1].Err)
	})
}
