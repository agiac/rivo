package core_test

import (
	"context"
	"fmt"
	. "github.com/agiac/rivo/core"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleConnect() {
	ctx := context.Background()

	g := Of("Hello", "Hello", "Hello")

	capitalize := Map(func(ctx context.Context, i string) string {
		return strings.ToUpper(i)
	})

	lowercase := Map(func(ctx context.Context, i string) string {
		return strings.ToLower(i)
	})

	resA := make([]string, 0)
	a := Do(func(ctx context.Context, i string) {
		resA = append(resA, i)
	})

	resB := make([]string, 0)
	b := Do(func(ctx context.Context, i string) {
		resB = append(resB, i)
	})

	p1 := Pipe(capitalize, a)
	p2 := Pipe(lowercase, b)

	<-Pipe(g, Connect(p1, p2))(ctx, nil)

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

		g := Of("Hello", "Hello", "Hello")

		mu1 := sync.Mutex{}
		resA := make([]string, 0)
		a := Do(func(ctx context.Context, i string) {
			mu1.Lock()
			resA = append(resA, i)
			mu1.Unlock()
		})

		mu2 := sync.Mutex{}
		resB := make([]string, 0)
		b := Do(func(ctx context.Context, i string) {
			mu2.Lock()
			resB = append(resB, i)
			mu2.Unlock()
		})

		<-Pipe(g, Connect(a, b))(ctx, nil)

		want := []string{"Hello", "Hello", "Hello"}

		mu1.Lock()
		mu2.Lock()
		defer mu1.Unlock()
		defer mu2.Unlock()

		assert.Equal(t, want, resA)
		assert.Equal(t, want, resB)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		g := Of("Hello", "Hello", "Hello")

		mtxA := sync.Mutex{}
		resA := make([]string, 0)
		a := Do(func(ctx context.Context, i string) {
			mtxA.Lock()
			defer mtxA.Unlock()
			resA = append(resA, i)
		})

		mtxB := sync.Mutex{}
		resB := make([]string, 0)
		b := Do(func(ctx context.Context, i string) {
			mtxB.Lock()
			defer mtxB.Unlock()
			resB = append(resB, i)
		})

		<-Connect(a, b)(ctx, g(ctx, nil))

		mtxA.Lock()
		mtxB.Lock()
		defer mtxA.Unlock()
		defer mtxB.Unlock()

		assert.LessOrEqual(t, len(resA), 3)
		assert.LessOrEqual(t, len(resB), 3)
	})
}
