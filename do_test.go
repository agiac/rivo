package rivo_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
)

func ExampleDo() {
	ctx := context.Background()

	in := make(chan rivo.Item[int])
	go func() {
		defer close(in)
		in <- rivo.Item[int]{Val: 1}
		in <- rivo.Item[int]{Val: 2}
		in <- rivo.Item[int]{Err: errors.New("error 1")}
		in <- rivo.Item[int]{Val: 4}
		in <- rivo.Item[int]{Err: errors.New("error 2")}
	}()

	d := rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
		if i.Err != nil {
			fmt.Printf("ERROR: %v\n", i.Err)
		}
	})

	<-d(ctx, in)

	// Output:
	// ERROR: error 1
	// ERROR: error 2
}

func TestDo(t *testing.T) {
	t.Run("do all items", func(t *testing.T) {
		ctx := context.Background()

		count := 0

		g := rivo.Of(1, 2, 3, 4, 5)

		d := rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
			count++
		})

		p := rivo.Pipe(g, d)

		<-p(ctx, nil)

		assert.Equal(t, 5, count)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		count := 0

		g := rivo.Of(1, 2, 3, 4, 5)

		d := rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
			count++
			if i.Val == 3 {
				cancel()
			}
		})

		p := rivo.Pipe(g, d)

		<-p(ctx, nil)

		assert.LessOrEqual(t, 4, count)
	})

	t.Run("with pool size", func(t *testing.T) {
		ctx := context.Background()

		count := atomic.Int32{}

		g := rivo.Of[int32](1, 2, 3, 4, 5)

		d := rivo.Do(func(ctx context.Context, i rivo.Item[int32]) {
			count.Add(1)
		}, rivo.WithPoolSize(3))

		p := rivo.Pipe(g, d)

		<-p(ctx, nil)

		assert.Equal(t, int32(5), count.Load())
	})
}
