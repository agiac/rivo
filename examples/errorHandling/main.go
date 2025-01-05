package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/agiac/rivo"
)

// This example demonstrates how to handle errors in a pipeline.
// The pipeline takes a sequence of strings, converts them to integers, doubles the values, and logs them.
// If a string cannot be converted to an integer, the error is logged in a concurrent pipeline.

func main() {
	ctx := context.Background()

	g := rivo.Of("1", "2", "3_", "4", "5", "6**", "?", "8", "9", "10")

	toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
		return strconv.Atoi(i.Val)
	})

	double := rivo.Map(func(ctx context.Context, i rivo.Item[int]) (int, error) {
		if i.Err != nil {
			return 0, i.Err // Pass errors along
		}
		return i.Val * 2, nil
	})

	logValue := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
		fmt.Printf("Value: %d\n", i.Val)
	})

	filterNonError := rivo.Filter[int](func(ctx context.Context, i rivo.Item[int]) (bool, error) {
		return i.Err == nil, nil
	})

	filterError := rivo.Filter[int](func(ctx context.Context, i rivo.Item[int]) (bool, error) {
		return i.Err != nil, nil
	})

	errs := make([]error, 0)
	handleError := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
		errs = append(errs, i.Err)
	})

	valP := rivo.Pipe3(filterNonError, double, logValue)

	errP := rivo.Pipe(filterError, handleError)

	<-rivo.Pipe4(g, toInt, double, rivo.Connect(valP, errP))(ctx, nil)

	for _, err := range errs {
		fmt.Printf("Error: %v\n", err)
	}

	// Expected output:
	// Value: 4
	// Value: 8
	// Value: 16
	// Value: 20
	// Value: 32
	// Value: 36
	// Value: 40
	// Error: strconv.Atoi: parsing "3_": invalid syntax
	// Error: strconv.Atoi: parsing "6**": invalid syntax
	// Error: strconv.Atoi: parsing "?": invalid syntax
}
