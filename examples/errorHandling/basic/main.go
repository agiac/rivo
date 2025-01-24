package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/agiac/rivo"
)

// This example demonstrates the simplest way to handle errors in a pipeline,
// by doing so together with the other values in one of more stages.

func main() {
	ctx := context.Background()

	g := rivo.Of("1", "2", "3_", "4", "5", "6**", "?", "8", "9", "10")

	toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
		if i.Err != nil {
			return 0, i.Err // Pass errors along
		}

		return strconv.Atoi(i.Val)
	})

	double := rivo.Map(func(ctx context.Context, i rivo.Item[int]) (int, error) {
		if i.Err != nil {
			return 0, i.Err // Pass errors along
		}

		return i.Val * 2, nil
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
