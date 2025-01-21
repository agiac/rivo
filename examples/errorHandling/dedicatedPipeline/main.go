package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/agiac/rivo"
)

// This example demonstrates the recommended approach to handle errors, in a separate dedicate pipeline.

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

	handleErrors := rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
		fmt.Printf("Error: %v\n", i.Err)
	})

	handleValues := rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
		fmt.Printf("Value: %d\n", i.Val)
	})

	errs, vals := rivo.SegregateErrors(rivo.Pipe3(g, toInt, double))(ctx, nil)

	<-rivo.Connect(rivo.Pipe(errs, handleErrors), rivo.Pipe(vals, handleValues))(ctx, nil)

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
