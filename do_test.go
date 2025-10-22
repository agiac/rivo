package rivo_test

import (
	"context"
	"fmt"
	. "github.com/agiac/rivo"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleDo() {
	ctx := context.Background()

	g := Of(1, 2, 3, 4, 5)

	d := Do(func(ctx context.Context, i int) error {
		fmt.Println(i)
		return nil
	})

	<-Pipe(g, d)(ctx, nil, nil)

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
		d := Do(func(ctx context.Context, i int) error {
			count++
			return nil
		})

		p := Pipe(g, d)

		<-p(ctx, nil, nil)

		assert.Equal(t, 5, count)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		count := 0

		g := Of(1, 2, 3, 4, 5)
		d := Do(func(ctx context.Context, i int) error {
			count++
			return nil
		})

		p := Pipe(g, d)

		<-p(ctx, nil, nil)

		assert.Lessf(t, count, 5, "count should be less than 5 when context is cancelled")
	})

	t.Run("with pool size", func(t *testing.T) {
		ctx := context.Background()

		count := atomic.Int32{}
		g := Of[int32](1, 2, 3, 4, 5)
		d := Do(func(ctx context.Context, i int32) error {
			count.Add(1)
			return nil
		}, DoPoolSize(3))

		p := Pipe(g, d)

		<-p(ctx, nil, nil)

		assert.Equal(t, int32(5), count.Load())
	})

	t.Run("with onBeforeClose", func(t *testing.T) {
		ctx := context.Background()

		count := atomic.Int32{}

		g := Of[int32](1, 2, 3, 4, 5)

		onBeforeCloseCalled := false

		d := Do(func(ctx context.Context, i int32) error {
			count.Add(1)
			return nil
		}, DoOnBeforeClose(func(ctx context.Context) {
			onBeforeCloseCalled = true
		}))

		p := Pipe(g, d)

		<-p(ctx, nil, nil)

		assert.Equal(t, int32(5), count.Load())
		assert.True(t, onBeforeCloseCalled)
	})

	t.Run("handle errors", func(t *testing.T) {
		ctx := context.Background()

		count := 0

		g := Of(1, 2, 3, 4, 5)
		d := Do(func(ctx context.Context, i int) error {
			if i == 3 {
				return fmt.Errorf("error on 3")
			}
			count++
			return nil
		})

		p := Pipe(g, d)

		var mu sync.Mutex
		var foundErrs []error
		errs, stop := RunErrorSync(ctx, func(ctx context.Context, errs <-chan error) {
			for err := range errs {
				mu.Lock()
				foundErrs = append(foundErrs, err)
				mu.Unlock()
			}
		})
		defer stop()

		<-p(ctx, nil, errs)

		assert.Equal(t, 4, count)
		mu.Lock()
		assert.Len(t, foundErrs, 1)
		assert.EqualError(t, foundErrs[0], "error on 3")
		mu.Unlock()
	})
}
