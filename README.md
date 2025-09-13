# rivo

[![Go Reference](https://pkg.go.dev/badge/github.com/agiac/rivo.svg)](https://pkg.go.dev/github.com/agiac/rivo)

`rivo` is a concurrent stream processing library for Go that provides type safety through generics and a composable pipeline architecture.

**NOTE: THIS LIBRARY IS STILL IN ACTIVE DEVELOPMENT AND IS NOT YET PRODUCTION READY.**

TODO: update for new Pipeline function signature

## About

`rivo` has two major inspirations:
1. The book ["Concurrency in Go"](https://www.amazon.com/Concurrency-Go-Tools-Techniques-Developers/dp/1491941197);
2. [ReactiveX](https://reactivex.io/), in particular the [Go](https://github.com/ReactiveX/RxGo) and [JS](https://github.com/ReactiveX/rxjs) libraries;

Compared to these sources, `rivo` aims to provide better type safety (both "Concurrency in Go" and RxGo were written in a pre-generics era and make heavy use of `interface{}`) 
and a more intuitive API and developer experience (Rx is very powerful, but can be overwhelming for newcomers).

## Getting started

### Prerequisites

`rivo` requires Go 1.23 or later. 

### Installation

```shell
  go get github.com/agiac/rivo
```

### Basic concepts

`rivo` has several main types, which are the building blocks of the library: `Stream`, `Pipeline`, `Generator`, `Sync`, and `Item`.

`Stream` represents a data stream. It is a read-only channel of type T.

```go
type Stream[T any] <-chan T
```

`Pipeline` is a function that takes a `context.Context` and a `Stream` of one type and returns a `Stream` of the same or a different type.
They represent the operations that can be performed on streams. Pipelines can be composed together to create more complex operations.

```go
type Pipeline[T, U any] func(ctx context.Context, stream Stream[T]) Stream[U]
```

For convenience, `rivo` also provides type aliases for common pipeline patterns:

```go
// Generator is a pipeline that generates items of type T without any input
type Generator[T any] = Pipeline[None, T]

// Sync is a pipeline that processes items of type T and does not emit any items
type Sync[T any] = Pipeline[T, None]
```

`Item` is a struct that contains a value and an optional error. It's used when you need error handling in your streams:

```go
type Item[T any] struct {
	Val T
	Err error
}
```

Most basic operations work with plain values, but when you need error handling, you can use `Item[T]` and the corresponding pipelines that support error propagation.

If a pipeline generates values without depending on an input stream, it is called a _generator_. 
If it consumes values without generating a new stream, it is called a _sink_. 
If it transforms values, it is called a _transformer_.

Here's a basic example:

```go
package main

import (
	"context"
	"fmt"
	"github.com/agiac/rivo"
)

func main() {
	ctx := context.Background()

	// `Of` returns a generator that returns a stream that will emit the provided values
	in := rivo.Of(1, 2, 3, 4, 5)

	// `Filter` returns a pipeline that filters the input stream using the given function
	onlyEven := rivo.Filter(func(ctx context.Context, n int) bool {
		return n%2 == 0
	})

    // `Do` returns a pipeline that applies the given function to each item in the input stream without emitting any values
	log := rivo.Do(func(ctx context.Context, n int) {
		fmt.Println(n)
	})

	// `Pipe` composes pipelines together, returning a new pipeline
	p := rivo.Pipe3(in, onlyEven, log)

	// By passing a context and an input channel to our pipeline, we can get the output stream.
	// Since our first pipeline `in` is a generator and does not depend on an input stream, we can pass a nil channel.
	// Also, since `log` is a sink, we only have to read once from the output channel to know that the pipeline has finished.
	<-p(ctx, nil)

	// Expected output:
	// 2
	// 4
}
```

For error handling scenarios, you can use `Item[T]` as your data type to carry both values and errors:

```go
package main

import (
	"context"
	"fmt"
	"strconv"
	"github.com/agiac/rivo"
)

func main() {
	ctx := context.Background()

	// Create a generator with string values
	g := rivo.Of("1", "2", "invalid", "4", "5")

	// Transform string to Item[int] with error handling
	toInt := rivo.Map(func(ctx context.Context, s string) rivo.Item[int] {
		n, err := strconv.Atoi(s)
		if err != nil {
			return rivo.Item[int]{Err: err} // Return an item with the error
		}
		return rivo.Item[int]{Val: n} // Return an item with the value
	})

	// Process the items, handling both values and errors
	handleResults := rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
		if i.Err != nil {
			fmt.Printf("ERROR: %v\n", i.Err)
		} else {
			fmt.Printf("Value: %d\n", i.Val)
		}
	})

	p := rivo.Pipe3(g, toInt, handleResults)
	<-p(ctx, nil)
	
	// Output:
	// Value: 1
	// Value: 2
	// ERROR: strconv.Atoi: parsing "invalid": invalid syntax
	// Value: 4
	// Value: 5
}
```

## Pipeline factories

`rivo` comes with a set of built-in pipeline factories.

### Generators
- `Of`: returns a generator pipeline that emits the provided values;
- `FromFunc`: returns a generator pipeline that emits values returned by the provided function until the function returns false;
- `FromSeq` and `FromSeq2`: return generator pipelines that emit the values from the provided iterators;
- `Tee` and `TeeN`: return N generator pipelines that each receive a copy of each item from the input stream;
- `Segregate`: returns two generator pipelines, where the first pipeline emits items that pass the predicate, and the second pipeline emits items that do not pass the predicate;

### Sinks
- `Do`: returns a sink pipeline that performs a side effect for each item in the input stream;
- `Connect`: returns a sink pipeline that applies the given sink pipelines to the input stream concurrently;

### Transformers
- `Filter`: returns a transformer pipeline that filters the input stream using the given function;
- `Map`: returns a transformer pipeline that applies a function to each item from the input stream;
- `FilterMap`: returns a transformer pipeline that filters and maps items from the input stream in a single operation;
- `Batch`: returns a transformer pipeline that groups the input stream into batches of the provided size;
- `Flatten`: returns a transformer pipeline that flattens the input stream of slices;
- `ForEachOutput`: returns a transformer pipeline that applies a function to each item, allowing direct output channel access;
- `Pipe`, `Pipe2`, `Pipe3`, `Pipe4`, `Pipe5`: return transformer pipelines that compose the provided pipelines together;

Besides these, the library's subdirectories contain more specialized pipeline factories.

### Package `rivo/io`

- `FromReader`: returns a generator pipeline that reads from the provided `io.Reader` and emits the read bytes;
- `ToWriter`: returns a sink pipeline that writes the input stream to the provided `io.Writer`;

### Package `rivo/bufio`

- `FromScanner`: returns a generator pipeline that reads from the provided `bufio.Scanner` and emits the scanned items;
- `ToWriter`: returns a sink pipeline that writes the input stream to the provided `bufio.Writer`;

### Package `rivo/csv`

- `FromReader`: returns a generator pipeline that reads from the provided `csv.Reader` and emits the read records;
- `ToWriter`: returns a sink pipeline that writes the input stream to the provided `csv.Writer`;

### Package `rivo/aws/dynamodb`

- `Scan`: returns a generator pipeline that scans the provided DynamoDB table and emits the scan output responses;
- `ScanItems`: returns a generator pipeline that scans the provided DynamoDB table and emits the items from the scan output responses;
- `BatchWrite`: returns a transformer pipeline that writes the input stream to the provided DynamoDB table using the BatchWriteItem API;
- `BatchPutItems`: returns a transformer pipeline that writes the input stream to the provided DynamoDB table using the BatchWriteItem API, but only for PutItem operations;

## Configuration Options

Many pipelines support configuration options to customize their behavior:

- **Pool Size**: Control the number of concurrent goroutines (e.g., `MapPoolSize`, `FilterPoolSize`, `DoPoolSize`)
- **Buffer Size**: Control the internal channel buffer size (e.g., `MapBufferSize`, `BatchBufferSize`)
- **Time-based Options**: Control time-based behavior (e.g., `BatchMaxWait`)
- **Lifecycle Hooks**: Add hooks for cleanup or finalization (e.g., `FromFuncOnBeforeClose`)

Example usage:

```go
// Map with custom pool size and buffer size
mapper := rivo.Map(transformFunc, rivo.MapPoolSize(5), rivo.MapBufferSize(100))

// Batch with time-based batching
batcher := rivo.Batch(10, rivo.BatchMaxWait(100*time.Millisecond))
```

## Utilities

`rivo` provides several utility functions to work with streams:

- `Collect`: collects all items from a stream into a slice
- `CollectWithContext`: like `Collect` but respects context cancellation
- `OrDone`: utility function that propagates context cancellation to streams
- `FilterMapValues`: extracts only successful values from Item streams
- `FilterMapErrors`: extracts only errors from Item streams
- `Merge`: merges multiple streams into a single stream

## Error handling

When you need error handling in your streams, you can use the `Item[T]` type to carry both values and errors through your pipelines. This allows you to handle errors at any point in the pipeline without stopping the entire stream.

The library provides several utilities for working with error-carrying streams:

- `FilterMapValues`: extracts only successful values from Item streams, filtering out errors
- `FilterMapErrors`: extracts only errors from Item streams, filtering out successful values
- `Segregate`: splits any stream based on a predicate function

See `examples/errorHandling` for comprehensive examples of different error handling patterns.

## Examples

More examples can be found in the [examples](./examples) folder.

---

## Contributing

Contributions are welcome! If you have any ideas, suggestions or bug reports, please open an issue or a pull request.

## Roadmap

- [ ] Review docs, in particular where "pipeline" is used instead of "generator", "sink" or "transformer"
- [ ] Add more pipelines, also using the [RxJS list of operators](https://rxjs.dev/guide/operators) as a reference:
  - [x] FilterMap (combines filter and map operations)
  - [x] ForEachOutput (direct output channel access)
  - [ ] Tap (side effects without modifying the stream)
  - [ ] Time-based operators (throttle, debounce, etc.)
  - [ ] SQL-like operators (join, group by, etc.)
  - [ ] More AWS integrations
- [ ] Add more utilities:
  - [x] Merge (combine multiple streams)
  - [ ] Zip (combine streams element-wise)
  - [ ] Take/Skip operators
- [ ] Performance optimizations and benchmarking
- [ ] Add more examples and tutorials

## License

`rivo` is licensed under the MIT license. See the [LICENSE](./LICENSE) file for details.