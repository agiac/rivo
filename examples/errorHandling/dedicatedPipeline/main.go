package main

import (
	"context"
	"fmt"
	"github.com/agiac/rivo"
	"strconv"
)

// This example demonstrates the recommended approach to handle errors, in a separate dedicate pipeline.

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

	basePipeline := rivo.Pipe3(g, toInt, double)

	handleErrors := rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
		if i.Err == nil {
			return // Skip items that do not have an error
		}
		fmt.Printf("Error: %v\n", i.Err)
	})

	handleValues := rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
		if i.Err != nil {
			return // Skip items that have an error
		}
		fmt.Printf("Value: %d\n", i.Val)
	})

	<-rivo.Connect(
		handleValues,
		handleErrors,
	)(ctx, basePipeline(ctx, nil))

	// Expected output (the order might be different because the handleErrors and handleValues pipeline run concurrently):
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
