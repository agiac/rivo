# rivo

`rivo` is a library for stream processing in Go.

**NOTE: THIS LIBRARY IS STILL IN ACTIVE DEVELOPMENT AND IS NOT YET PRODUCTION READY.**

## About

`rivo` had two major inspirations:
1. The book ["Concurrency in Go"](https://www.amazon.com/Concurrency-Go-Tools-Techniques-Developers/dp/1491941197);
2. [ReactiveX](https://reactivex.io/), in particular the [Go](https://github.com/ReactiveX/RxGo) and [JS](https://github.com/ReactiveX/rxjs) libraries;

Compared to these sources, `rivo` aims to provide better type safety (both "Concurrency in Go" and RxGo were written in a pre-generics era and make a heavy use of `interface{}`) 
and a more intuitive API and developer experience (Rx is very powerful, but can be overwhelming for newcomers).

## Getting started

### Prerequisites

`rivo` requires Go 1.23 or later. 

 ### Installation

```shell
  go get github.com/agiac/rivo
```

### Basic concepts

`rivo` has 3 main types, which are the building blocks of the library: `Item`, `Stream` and `Pipeline`.

`Item` is a struct which contains a value and an optional error. Just like errors are returned next to the result
of a function in synchronous code, they should be passed along into asynchronous code and handled where more appropriate.

```go
type Item[T any] struct {
	Val T
	Err error
}
```

`Stream` is a read only channel of items. As the name suggests, it represents a stream of data.

```go
type Stream[T any] <-chan Item[T]
```

`Pipeline` is a function that takes a `context.Context` and a `Stream` of one type and returns a `Stream` of the same or a different type.
They represent the operations that can be performed on streams. Pipelines can be composed together to create more complex operations.

```go
type Pipeline[T, U any] func(ctx context.Context, stream Stream[T]) Stream[U]
```

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

	// `Of` returns a generator which returns a stream that will emit the provided values
	in := rivo.Of(1, 2, 3, 4, 5)

	// `Filter` returns a transformer that filters the input stream using the given function.
	onlyEven := rivo.Filter(func(ctx context.Context, i rivo.Item[int]) (bool, error) {
		// Always check for errors
		if i.Err != nil {
			return true, i.Err // Propagate the error
		}

		return i.Val%2 == 0, nil
	})

    // `Do` returns a sync that applies the given function to each item in the input stream, without emitting any values.
	log := rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
		if i.Err != nil {
			fmt.Printf("ERROR: %v\n", i.Err)
			return
		}

		fmt.Println(i.Val)
	})

	// `Pipe` composes pipelines together, returning a new pipeline.
	p := rivo.Pipe3(in, onlyEven, log)

	// By passing a context and an input channel to our pipeline, we can get the output stream.
	// Since our first pipeline `in` is a generator and does not depend on an input stream, we can pass a nil channel.
	// Also, since log is a sink, we only have to read once from the output channel to know that the pipe has finished.
	<-p(ctx, nil)

	// Expected output:
	// 2
	// 4
}
```

## Pipeline factories

`rivo` comes with a set of built-in pipeline factories.

### Generators
- `Of`: returns a pipeline which returns a stream that will emit the provided values;
- `FromFunc`: returns a pipeline which returns a stream that will emit the values returned by the provided function;
- `FromSeq` and `FromSeq2`: returns a pipeline which returns a stream that will emit the values from the provided iterator;

### Sinks
- `Do`: returns a pipeline which performs a side effect for each item in the input stream;

### Transformers
- `Filter`: returns a pipeline which filters the input stream using the given function;
- `Map`: returns a pipeline which maps the input stream using the given function;
- `ForEach`: returns a pipeline which applies the given function to each item in the input stream and forwards only the errors;
- `Batch`: returns a pipeline which groups the input stream into batches of the provided size;
- `Flatten`: returns a pipeline which flattens the input stream of slices; 
- `Pipe`, `Pipe2`, `Pipe3`, `Pipe4`, `Pipe5`: return a pipeline which composes the provided pipelines together;
- `Connect`: returns a sync which applies the given syncs to the input stream concurrently;
- `Segregate`: returns a function that returns two pipelines, where the first pipeline emits items that pass the predicate, and the second pipeline emits items that do not pass the predicate.
- `Tee` and `TeeN`: return n pipelines that each receive a copy of each item from the input stream;

Besides these, the directories of the library contain more specialized pipelines factories.

### Package `rivo/io`

- `FromReader`: returns a pipeline which reads from the provided `io.Reader` and emits the read bytes;
- `ToWriter`: returns a pipeline which writes the input stream to the provided `io.Writer`;

### Package `rivo/bufio`

- `FromScanner`: returns a pipeline which reads from the provided `bufio.Scanner` and emits the scanned items;
- `ToScanner`: returns a pipeline which writes the input stream to the provided `bufio.Writer`;

### Package `rivo/csv`

- `FromReader`: returns a pipeline which reads from the provided `csv.Reader` and emits the read records;
- `ToWriter`: returns a pipeline which writes the input stream to the provided `csv.Writer`;

### Package `rivio/errors`

- `WithErrorHandler`: returns a pipeline that connects the input pipeline to an error handling pipeline.

## Optional parameters

Many pipeline factories accepts a common set of optional parameters. These can be provided via functional options.

```go
  double := rivo.Map(
	  func(ctx context.Context, i rivo.Item[int]) (int, error) { return i.Val * 2, nil  },
	  // `Pass additional options to the pipeline
	  rivo.WithBufferSize(1), 
	  rivo.WithPoolSize(runtime.NumCPU()), 
	  )
```

The currently available options are:

- `WithPoolSize(int)`: sets the number of goroutines that will be used to process items. Default is 1.
- `WithBufferSize(int)`: sets the buffer size of the output channel. Default is 0;
- `WithStopOnError(bool)`: if true, the pipeline will stop processing items when an error is encountered. Default is false.
- `WithOnBeforeClosed(func(context.Context) error)`: a function that will be called before the output channel is closed.




## Error handling

As mentioned, each values contains a value and an optional error. You can handle error either individually inside pipelines' callbacks like `Map` or `Do` or
(recommended) create dedicated pipelines for error handling. See `examples/errorHanlidng` for this regard.

## Examples

More examples can be found in the [examples](./examples) folder.

---

## Contributing

Contributions are welcome! If you have any ideas, suggestions or bug reports, please open an issue or a pull request.

## Roadmap

- [ ] Review docs, in particular where "pipeline" is used instead of "generator", "sink" or "transformer"
- [ ] Remove current dedicated folders for special pipelines and move them to the main package
- [ ] Consider dedicated options for each pipeline instead of a common set of options
- [ ] Add more pipelines, also using the [RxJS list of operators](https://rxjs.dev/guide/operators) as a reference:
  - [ ] Tap 
  - [ ] Better error handling
  - [ ] Time-based
  - [ ] SQL
  - [ ] AWS
  - [ ] ...
- [ ] Add more utilities
  - [ ] Merge
- [ ] Add more examples
- [ ] Error handling section in the README

## License

`rivo` is licensed under the MIT license. See the [LICENSE](./LICENSE) file for details.









