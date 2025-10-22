package main

import (
	"context"
	"fmt"

	"github.com/agiac/rivo"
)

// This example demonstrates how to pass additional options to a pipeline.

func main() {
	ctx := context.Background()

	in := rivo.Of(1, 2, 3, 4, 5)

	doubleFn := func(ctx context.Context, n int) (int, error) {
		return n * 2, nil
	}

	// `Pass additional options to the pipeline
	double := rivo.Map(doubleFn, rivo.MapBufferSize(1))

	p := rivo.Pipe(in, double)

	s := p(ctx, nil, nil)

	for n := range s {
		fmt.Println(n)
	}

	// Output:
	// 2
	// 4
	// 6
	// 8
	// 10
}
