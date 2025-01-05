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

	// `Of` returns a generator which returns a stream that will emit the provided values
	in := rivo.Of(1, 2, 3, 4, 5)

	// `Filter` returns a pipeable that filters the input stream using the given function.
	onlyEven := rivo.Filter(func(ctx context.Context, i rivo.Item[int]) (bool, error) {
		// Always check for errors
		if i.Err != nil {
			return true, i.Err // Propagate the error
		}

		return i.Val%2 == 0, nil
	})

	log := rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
		if i.Err != nil {
			fmt.Printf("ERROR: %v\n", i.Err)
			return
		}

		fmt.Println(i.Val)
	})

	// `Pipe` composes pipeables together, returning a new pipeable
	p := rivo.Pipe3(in, onlyEven, log)

	// By passing a context and an input channel to our pipeable, we can get the output stream.
	// Since our first pipeable `in` is a generator and does not depend on an input stream, we can pass a nil channel.
	// Also, since log is a sink, we only have to read once from the output channel to know that the pipe has finished.
	<-p(ctx, nil)

	// Expected output:
	// 2
	// 4
}
