package rivo_test

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

func ExampleWithErrorHandler() {
	ctx := context.Background()

	g := rivo.Of("1", "_2", "*", "4", "5")

	toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
		return strconv.Atoi(i.Val)
	})

	handleErrors := rivo.Do(func(ctx context.Context, i rivo.Item[error]) {
		fmt.Printf("Error: %v\n", i.Val)
	})

	vals := rivo.Collect(rivo.WithErrorHandler(rivo.Pipe(g, toInt), handleErrors)(ctx, nil))

	for _, v := range vals {
		fmt.Printf("Value: %d\n", v.Val)
	}

	// Output:
	// Error: strconv.Atoi: parsing "_2": invalid syntax
	// Error: strconv.Atoi: parsing "*": invalid syntax
	// Value: 1
	// Value: 4
	// Value: 5
}

func TestWithErrorHandler(t *testing.T) {
	t.Run("one pipeline", func(t *testing.T) {
		ctx := context.Background()

		g := rivo.Of("1", "_2", "*", "4", "5")

		toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
			return strconv.Atoi(i.Val)
		})

		errs := make([]error, 0)
		handleErrors := rivo.Do(func(ctx context.Context, i rivo.Item[error]) {
			errs = append(errs, i.Val)
		})

		vals := rivo.Collect(rivo.WithErrorHandler(rivo.Pipe(g, toInt), handleErrors)(ctx, nil))

		assert.Equal(t, []rivo.Item[int]{{Val: 1}, {Val: 4}, {Val: 5}}, vals)
		assert.Equal(t, 2, len(errs))
	})

	t.Run("multiple pipelines", func(t *testing.T) {
		ctx := context.Background()

		g := rivo.Of("1", "_2", "*", "4", "5")

		toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
			return strconv.Atoi(i.Val)
		})

		double := rivo.Map(func(ctx context.Context, i rivo.Item[int]) (int, error) {
			if i.Val == 4 {
				return 0, assert.AnError
			}

			return i.Val * 2, nil
		})

		mtx := sync.Mutex{}
		errs := make([]error, 0)
		handleErrors := rivo.Do(func(ctx context.Context, i rivo.Item[error]) {
			mtx.Lock()
			defer mtx.Unlock()
			errs = append(errs, i.Val)
		})

		vals := rivo.Collect(rivo.Pipe(
			rivo.WithErrorHandler(rivo.Pipe(g, toInt), handleErrors),
			rivo.WithErrorHandler(double, handleErrors),
		)(ctx, nil))

		assert.Equal(t, []rivo.Item[int]{{Val: 2}, {Val: 10}}, vals)
		assert.Equal(t, 3, len(errs))
	})
}
