package main

import (
	"context"
	"fmt"
	rivo "github.com/agiac/rivo/core"
	"strconv"
)

// This example demonstrates the simplest way to handle errors in a pipeline:
// we use rivo.Item to pass both values and errors through the pipeline,
// allowing us to handle errors at any point in the pipeline without stopping the entire stream.
// In this example, the errors are propagated through the pipeline, until we reach the end where we handle them,
// together with the values, in a single place.

func main() {
	ctx := context.Background()

	g := rivo.Of("1", "2", "3_", "4", "5", "6**", "?", "8", "9", "10")

	toInt := rivo.Map(func(ctx context.Context, s string) rivo.Item[int] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return rivo.Item[int]{Err: err} // Return an item with the error
		}
		return rivo.Item[int]{Val: n} // Return an item with the value
	})

	double := rivo.Map(func(ctx context.Context, i rivo.Item[int]) rivo.Item[int] {
		if i.Err != nil {
			return i // If there's an error, return it as is
		}

		return rivo.Item[int]{Val: i.Val * 2} // Otherwise, double the value
	})

	handleValuesAndErrors := rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
		if i.Err != nil {
			fmt.Printf("Error: %v\n", i.Err)
		} else {
			fmt.Printf("Value: %d\n", i.Val)
		}
	})

	<-rivo.Pipe4(g, toInt, double, handleValuesAndErrors)(ctx, nil)

	// Expected output:
	// Value: 2
	// Value: 4
	// Error: strconv.Atoi: parsing "3_": invalid syntax
	// Value: 8
	// Value: 10
	// Error: strconv.Atoi: parsing "6**": invalid syntax
	// Error: strconv.Atoi: parsing "?": invalid syntax
	// Value: 16
	// Value: 18
	// Value: 20
}
