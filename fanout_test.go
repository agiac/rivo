package rivo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleFanout() {
	ctx := context.Background()

	g := Of("Hello", "Hello", "Hello")

	capitalize := Map(func(ctx context.Context, i string) (string, error) {
		return strings.ToUpper(i), nil
	})

	lowercase := Map(func(ctx context.Context, i string) (string, error) {
		return strings.ToLower(i), nil
	})

	resA := make([]string, 0)
	a := Do(func(ctx context.Context, i string) error {
		resA = append(resA, i)
		return nil
	})

	resB := make([]string, 0)
	b := Do(func(ctx context.Context, i string) error {
		resB = append(resB, i)
		return nil
	})

	p1 := Pipe(capitalize, a)
	p2 := Pipe(lowercase, b)

	<-Fanout(g, p1, p2)(ctx, nil, nil)

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

func TestFanout(t *testing.T) {
	t.Run("fans out to all downstream workers", func(t *testing.T) {
		ctx := context.Background()

		g := Of("a", "b", "c")

		muA := sync.Mutex{}
		resA := make([]string, 0)
		a := Do(func(ctx context.Context, i string) error {
			muA.Lock()
			resA = append(resA, strings.ToUpper(i))
			muA.Unlock()
			return nil
		})

		muB := sync.Mutex{}
		resB := make([]string, 0)
		b := Do(func(ctx context.Context, i string) error {
			muB.Lock()
			resB = append(resB, strings.ToLower(i))
			muB.Unlock()
			return nil
		})

		<-Fanout(g, a, b)(ctx, nil, nil)

		muA.Lock()
		muB.Lock()
		defer muA.Unlock()
		defer muB.Unlock()

		assert.Equal(t, []string{"A", "B", "C"}, resA)
		assert.Equal(t, []string{"a", "b", "c"}, resB)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		g := Of(1, 2, 3, 4, 5)

		mu := sync.Mutex{}
		n := 0
		a := Do(func(ctx context.Context, i int) error {
			time.Sleep(2 * time.Millisecond)
			mu.Lock()
			n++
			mu.Unlock()
			return nil
		})

		<-Fanout(g, a)(ctx, nil, nil)

		mu.Lock()
		defer mu.Unlock()
		assert.LessOrEqual(t, n, 1, "expected few or no items due to context cancellation")
	})

	t.Run("propagates errors to errs channel", func(t *testing.T) {
		ctx := context.Background()

		g := Of(1, 2, 3)

		errBoom := errors.New("boom")
		bad := Do(func(ctx context.Context, i int) error {
			if i == 2 {
				return errBoom
			}
			return nil
		})

		errs := make(chan error, 10)
		<-Fanout(g, bad)(ctx, nil, errs)
		close(errs)

		var got []error
		for err := range errs {
			got = append(got, err)
		}

		assert.Contains(t, got, errBoom)
	})
}
