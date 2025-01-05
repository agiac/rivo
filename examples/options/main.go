package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/agiac/rivo"
)

// This example demonstrates how to pass additional options to a pipeable.

func main() {
	ctx := context.Background()

	in := rivo.Of(1, 2, 3, 4, 5)

	doubleFn := func(ctx context.Context, i rivo.Item[int]) (int, error) {
		if i.Err != nil {
			return 0, i.Err
		}

		// Simulate an error
		if i.Val == 3 {
			return 0, errors.New("some error")
		}

		return i.Val * 2, nil
	}

	// `Pass additional options to the pipeable
	double := rivo.Map(doubleFn, rivo.WithBufferSize(1), rivo.WithStopOnError(true))

	p := rivo.Pipe(in, double)

	s := p(ctx, nil)

	for item := range s {
		if item.Err != nil {
			fmt.Printf("ERROR: %v\n", item.Err)
			continue
		}
		fmt.Println(item.Val)
	}

	// Output:
	// 2
	// 4
	// ERROR: some error
}
