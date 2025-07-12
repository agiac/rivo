package rivo_test

import (
	"context"
	"fmt"
	. "github.com/agiac/rivo"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleDo() {
	ctx := context.Background()

	g := Of(1, 2, 3, 4, 5)

	d := Do(func(ctx context.Context, i int) {
		fmt.Println(i)
	})

	<-Pipe(g, d)(ctx, nil)

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func TestDo(t *testing.T) {
	t.Run("do all items", func(t *testing.T) {
		ctx := context.Background()

		count := 0

		g := Of(1, 2, 3, 4, 5)

		d := Do(func(ctx context.Context, i int) {
			count++
		})

		p := Pipe(g, d)

		<-p(ctx, nil)

		assert.Equal(t, 5, count)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		count := 0

		g := Of(1, 2, 3, 4, 5)

		d := Do(func(ctx context.Context, i int) {
			count++
		})

		p := Pipe(g, d)

		<-p(ctx, nil)

		assert.Lessf(t, count, 5, "count should be less than 5 when context is cancelled")
	})

	t.Run("with pool size", func(t *testing.T) {
		ctx := context.Background()

		count := atomic.Int32{}

		g := Of[int32](1, 2, 3, 4, 5)

		d := Do(func(ctx context.Context, i int32) {
			count.Add(1)
		}, DoPoolSize(3))

		p := Pipe(g, d)

		<-p(ctx, nil)

		assert.Equal(t, int32(5), count.Load())
	})

	t.Run("with onBeforeClose", func(t *testing.T) {
		ctx := context.Background()

		count := atomic.Int32{}

		g := Of[int32](1, 2, 3, 4, 5)

		onBeforeCloseCalled := false

		d := Do(func(ctx context.Context, i int32) {
			count.Add(1)
		}, DoOnBeforeClose(func(ctx context.Context) {
			onBeforeCloseCalled = true
		}))

		p := Pipe(g, d)

		<-p(ctx, nil)

		assert.Equal(t, int32(5), count.Load())
		assert.True(t, onBeforeCloseCalled)
	})
}
