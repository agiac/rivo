package main

import (
	"context"
	"fmt"
	"github.com/agiac/rivo"
)

// This example demonstrates a basic usage of pipeables and the Pipe function.
// We create a stream of integers and filter only the even ones.

func main() {
	ctx := context.Background()

	// `Of` is a factory function which returs a pipeable which returns a stream that will emit the provided values
	in := rivo.Of(1, 2, 3, 4, 5)

	// `Filter` returns a Pipeable that filters the input stream using the given function.
	onlyEven := rivo.Filter(func(ctx context.Context, i rivo.Item[int]) (bool, error) {
		// Always check for errors
		if i.Err != nil {
			return true, i.Err // Propagate the error
		}

		return i.Val%2 == 0, nil
	})

	// `Pipe` composes pipeables togheter, returnin a new pipeable
	p := rivo.Pipe(in, onlyEven)

	// By passing a context and an input channel to our pipeable, we can get the output stream.
	// Since our first pipeable `in` does not depend on a input stream, we can pass a nil channel.
	s := p(ctx, nil)

	// Consume the result stream
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
}
