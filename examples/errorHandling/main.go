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

	handleError := rivo.Do[int](func(ctx context.Context, i rivo.Item[int]) {
		fmt.Printf("Error: %v\n", i.Err)
	})

	valP := rivo.Pipe3(filterNonError, double, logValue)

	errP := rivo.Pipe(filterError, handleError)

	<-rivo.Pipe3(g, toInt, rivo.Connect(valP, errP))(ctx, nil)

	// Expected output:
	// Value: 2
	// Value: 4
	// Error: strconv.Atoi: parsing "3_": invalid syntax
	// Error: strconv.Atoi: parsing "6**": invalid syntax
	// Value: 8
	// Value: 10
	// Value: 16
	// Value: 18
	// Error: strconv.Atoi: parsing "?": invalid syntax
	// Value: 20
}
