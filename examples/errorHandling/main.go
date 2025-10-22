package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/agiac/rivo"
)

// This example demonstrates simple error handling in a pipeline.
// We create a stream of strings, convert them to integers, and log any conversion errors.

func main() {
	ctx := context.Background()

	// Create a generator with string values
	g := rivo.Of("1", "2", "invalid", "4", "5")

	// Transform string to Item[int] with error handling
	toInt := rivo.Map(func(ctx context.Context, s string) (int, error) {
		return strconv.Atoi(s)
	})

	// Handle successful integer values
	handleValues := rivo.Do[int](func(ctx context.Context, i int) error {
		fmt.Println("Value:", i)
		return nil
	})

	// Handle errors
	handleErrors := func(ctx context.Context, err error) {
		fmt.Println("ERROR:", err)
	}

	// Run the error handler in a separate goroutine
	errs, wait := rivo.RunErrorSyncFunc(ctx, handleErrors)
	defer wait()

	p := rivo.Pipe3(g, toInt, handleValues)

	<-p(ctx, nil, errs)

	// Value: 1
	// Value: 2
	// ERROR: strconv.Atoi: parsing "invalid": invalid syntax
	// Value: 4
	// Value: 5
}
