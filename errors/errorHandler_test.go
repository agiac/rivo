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

	g := rivo.Of("1", "2", "3_", "4", "5", "6**", "?", "8", "9", "10")

	toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
		return strconv.Atoi(i.Val)
	})

	double := rivo.Map(func(ctx context.Context, i rivo.Item[int]) (int, error) {
		return i.Val * 2, nil
	})

	logValue := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
		fmt.Printf("Value: %d\n", i.Val)
	})

	errs := make([]error, 0)
	handleError := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
		errs = append(errs, i.Err)
	})

	<-rivo.Pipe4(g, toInt, errors.WithErrorHandler(double, handleError), logValue)(ctx, nil)

	for _, err := range errs {
		fmt.Printf("Error: %v\n", err)
	}

	// Output:
	// Value: 2
	// Value: 4
	// Value: 8
	// Value: 10
	// Value: 16
	// Value: 18
	// Value: 20
	// Error: strconv.Atoi: parsing "3_": invalid syntax
	// Error: strconv.Atoi: parsing "6**": invalid syntax
	// Error: strconv.Atoi: parsing "?": invalid syntax
}

func TestWithErrorHandler(t *testing.T) {
	t.Run("with generator", func(t *testing.T) {
		ctx := context.Background()

		in := make(chan rivo.Item[int])
		go func() {
			defer close(in)
			in <- rivo.Item[int]{Val: 1}
			in <- rivo.Item[int]{Val: 2}
			in <- rivo.Item[int]{Err: fmt.Errorf("error")}
			in <- rivo.Item[int]{Val: 4}
			in <- rivo.Item[int]{Val: 5}
			in <- rivo.Item[int]{Err: fmt.Errorf("error")}
			in <- rivo.Item[int]{Err: fmt.Errorf("error")}
			in <- rivo.Item[int]{Val: 8}
			in <- rivo.Item[int]{Val: 9}
			in <- rivo.Item[int]{Val: 10}
		}()

		double := rivo.Map(func(ctx context.Context, i rivo.Item[int]) (int, error) {
			return i.Val * 2, nil
		})

		vals := make([]int, 0)
		handleValue := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
			vals = append(vals, i.Val)
		})

		errs := make([]error, 0)
		handleError := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
			errs = append(errs, i.Err)
		})

		<-rivo.Pipe2(errors.WithErrorHandler(double, handleError), handleValue)(ctx, in)

		expectedVals := []int{2, 4, 8, 10, 16, 18, 20}
		expectedErrs := []error{
			fmt.Errorf("error"),
			fmt.Errorf("error"),
			fmt.Errorf("error"),
		}

		assert.Equal(t, expectedVals, vals)
		assert.Equal(t, expectedErrs, errs)
	})

	t.Run("with transformer", func(t *testing.T) {
		ctx := context.Background()

		g := rivo.Of("1", "2", "3_", "4", "5", "6**", "?", "8", "9", "10")

		toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
			return strconv.Atoi(i.Val)
		})

		double := rivo.Map(func(ctx context.Context, i rivo.Item[int]) (int, error) {
			return i.Val * 2, nil
		})

		vals := make([]int, 0)
		handleValue := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
			vals = append(vals, i.Val)
		})

		errs := make([]error, 0)
		handleError := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
			errs = append(errs, i.Err)
		})

		<-rivo.Pipe4(g, toInt, errors.WithErrorHandler(double, handleError), handleValue)(ctx, nil)

		expectedVals := []int{2, 4, 8, 10, 16, 18, 20}
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

	t.Run("with sync", func(t *testing.T) {
		ctx := context.Background()

		g := rivo.Of("1", "2", "3_", "4", "5", "6**", "?", "8", "9", "10")

		toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
			return strconv.Atoi(i.Val)
		})

		double := rivo.Map(func(ctx context.Context, i rivo.Item[int]) (int, error) {
			return i.Val * 2, nil
		})

		vals := make([]int, 0)
		handleValue := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
			vals = append(vals, i.Val)
		})

		errs := make([]error, 0)
		handleError := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
			errs = append(errs, i.Err)
		})

		<-rivo.Pipe3(g, toInt, errors.WithErrorHandler(rivo.Pipe(double, handleValue), handleError))(ctx, nil)

		expectedVals := []int{2, 4, 8, 10, 16, 18, 20}
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

}

func TestWithErrorHandler2(t *testing.T) {
	t.Run("filter non errors and handle errors", func(t *testing.T) {
		ctx := context.Background()

		g := rivo.Of("1", "2", "3_", "4", "5", "6**", "?", "8", "9", "10")

		toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
			return strconv.Atoi(i.Val)
		})

		errs := make([]error, 0)
		handleError := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
			errs = append(errs, i.Err)
		})

		vals := make([]int, 0)
		handleValue := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
			vals = append(vals, i.Val)
		})

		<-rivo.Pipe3(g, errors.WithErrorHandler2(toInt, handleError), handleValue)(ctx, nil)

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
}
