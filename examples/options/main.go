package main

import (
	"context"
	"fmt"

	rivo "github.com/agiac/rivo/core"
)

// This example demonstrates how to pass additional options to a pipeline.

func main() {
	ctx := context.Background()

	in := rivo.Of(1, 2, 3, 4, 5)

	doubleFn := func(ctx context.Context, n int) int {
		return n * 2
	}

	// `Pass additional options to the pipeline
	double := rivo.Map(doubleFn, rivo.MapBufferSize(1))

	p := rivo.Pipe(rivo.Pipeline[rivo.None, int](in), double)

	s := p(ctx, nil)

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
