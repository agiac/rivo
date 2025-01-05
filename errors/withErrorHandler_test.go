package errors_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/agiac/rivo"
	"github.com/agiac/rivo/errors"
	"github.com/stretchr/testify/assert"
)

func ExampleWithErrorHandler() {
	ctx := context.Background()

	g := rivo.Of("1", "2", "3_", "4", "5**")

	toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
		return strconv.Atoi(i.Val)
	})

	logValue := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
		fmt.Printf("Value: %d\n", i.Val)
	})

	errs := make([]error, 0)
	errorHandler := rivo.Do[struct{}](func(ctx context.Context, i rivo.Item[struct{}]) {
		errs = append(errs, i.Err)
	})

	<-rivo.Pipe3(g, errors.WithErrorHandler(toInt, errorHandler), logValue)(ctx, nil)

	for _, err := range errs {
		fmt.Printf("Error: %v\n", err)
	}

	// Output:
	// Value: 1
	// Value: 2
	// Value: 4
	// Error: strconv.Atoi: parsing "3_": invalid syntax
	// Error: strconv.Atoi: parsing "5**": invalid syntax
}

func TestWithErrorHandler(t *testing.T) {
	t.Run("single source", func(t *testing.T) {
		ctx := context.Background()

		g := rivo.Of("1", "2", "3_", "4", "5", "6**", "?", "8", "9", "10")

		toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
			return strconv.Atoi(i.Val)
		})

		errs := make([]error, 0)
		errorHandler := rivo.Do(func(ctx context.Context, i rivo.Item[struct{}]) {
			errs = append(errs, i.Err)
		})

		vals := make([]int, 0)
		handleValue := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
			vals = append(vals, i.Val)
		})

		<-rivo.Pipe3(g, errors.WithErrorHandler(toInt, errorHandler), handleValue)(ctx, nil)

		expectedVals := []int{1, 2, 4, 5, 8, 9, 10}
		expectedErrs := []error{
			fmt.Errorf("strconv.Atoi: parsing \"3_\": invalid syntax"),
			fmt.Errorf("strconv.Atoi: parsing \"6**\": invalid syntax"),
			fmt.Errorf("strconv.Atoi: parsing \"?\": invalid syntax"),
		}

		assert.Equal(t, expectedVals, vals)
		for i, err := range errs {
			assert.EqualError(t, err, expectedErrs[i].Error())
		}
	})

	t.Run("multiple sources", func(t *testing.T) {
		ctx := context.Background()

		g1 := rivo.Of("1", "2", "3_", "4", "5")
		g2 := rivo.Of("true", "false", "_tt", "ff", "true")

		toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
			return strconv.Atoi(i.Val)
		})

		toBool := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (bool, error) {
			return strconv.ParseBool(i.Val)
		})

		errs := make([]error, 0)
		errorHandler := rivo.Do(func(ctx context.Context, i rivo.Item[struct{}]) {
			errs = append(errs, i.Err)
		})

		ints := make([]int, 0)
		handleInt := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
			ints = append(ints, i.Val)
		})

		bools := make([]bool, 0)
		handleBool := rivo.Do[bool](func(ctx context.Context, i rivo.Item[bool]) {
			bools = append(bools, i.Val)
		})

		<-rivo.Pipe3(g1, errors.WithErrorHandler(toInt, errorHandler), handleInt)(ctx, nil)
		<-rivo.Pipe3(g2, errors.WithErrorHandler(toBool, errorHandler), handleBool)(ctx, nil)

		expectedInts := []int{1, 2, 4, 5}
		expectedBools := []bool{true, false, true}
		expectedErrs := []error{
			fmt.Errorf("strconv.Atoi: parsing \"3_\": invalid syntax"),
			fmt.Errorf("strconv.ParseBool: parsing \"_tt\": invalid syntax"),
			fmt.Errorf("strconv.ParseBool: parsing \"ff\": invalid syntax"),
		}

		assert.Equal(t, expectedInts, ints)
		assert.Equal(t, expectedBools, bools)
		for i, err := range errs {
			assert.EqualError(t, err, expectedErrs[i].Error())
		}
	})
}
