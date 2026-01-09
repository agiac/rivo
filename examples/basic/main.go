package main

import (
	"context"
	"fmt"

	"github.com/agiac/rivo"
)

// This example demonstrates a basic usage of workers and the Pipe function.
// We create a stream of integers and filter only the even ones.

func main() {
	ctx := context.Background()

	// `Of` returns a generator which returns a stream that will emit the provided values
	in := rivo.Of(1, 2, 3, 4, 5)

	// `Filter` returns a worker that filters the input stream using the given function.
	onlyEven := rivo.Filter(func(ctx context.Context, n int) (bool, error) {
		return n%2 == 0, nil
	})

	// `Do` returns a worker that applies the given function to each item in the input stream, without emitting any values.
	log := rivo.Do(func(ctx context.Context, n int) error {
		fmt.Println(n)
		return nil
	})

	// `Pipe` composes workers together, returning a new worker
	p := rivo.Pipe3(in, onlyEven, log)

	// By passing a context and an input channel to our worker, we can get the output stream.
	// Since our first worker `in` is a generator and does not depend on an input stream, we can pass a nil channel.
	// Also, since log is a sink, we only have to read once from the output channel to know that the pipe has finished.
	<-p(ctx, nil, nil)

	// Expected output:
	// 2
	// 4
}
