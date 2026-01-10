# rivo

[![Go Reference](https://pkg.go.dev/badge/github.com/agiac/rivo.svg)](https://pkg.go.dev/github.com/agiac/rivo)

`rivo` is a library for highly concurrent Go programs that provides type safety through generics and a composable workers architecture.

**NOTE: THIS LIBRARY IS STILL IN ACTIVE DEVELOPMENT AND IS NOT YET PRODUCTION READY.**

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

`rivo` is built around its main `Worker` type, which has the following signature:

```go
type Worker[T, U any] = func(ctx context.Context, in <-chan T, errs chan<- error) <-chan U
```

This is the result of various iterations and refinements and I believe you could go a long way with it,
even without using the rest of the library.
  - The first argument is a `context.Context`, which allows for cancellation and timeouts, as well as passing values down the call chain, if needed. 
  - The second argument is a read-only channel of type `T`, which represents the input stream of data.
  - The third argument is a write-only channel for errors. Following Go's focus on explicit error handling, 
    this channel allows workers to report errors without stopping the entire stream.
  - The return value is a read-only channel of type `U`, which represents the output stream of data.

This structure allows for a clear flow of data, as well as composability of workers to create complex data processing pipelines, while maintaining type safety and explicit error handling.

For convenience, `rivo` also provides type aliases for common worker patterns:

`Generator` is a worker that generates items of type `T` without any input: 

```go
type Generator[T any] = Worker[None, T]
```

`Sync` is a worker that processes items of type `T` and does not emit any items, except after closing the output channel to signal completion:

```go
type Sync[T any] = Worker[T, None]
```

Here's a basic example:

```go
    package main
    
    import (
      "context"
      "fmt"
    
      "github.com/agiac/rivo"
    )
    
    // This example demonstrates a basic usage of workers and the Pipe function.
    // We create a channel of integers and filter only the even ones.
    
    func main() {
      ctx := context.Background()
    
      // `Of` returns a generator which returns a channel that will emit the provided values
      in := rivo.Of(1, 2, 3, 4, 5)
    
      // `Filter` returns a worker that filters the input channel using the given function.
      onlyEven := rivo.Filter(func(ctx context.Context, n int) (bool, error) {
        return n%2 == 0, nil
      })
    
      // `Do` returns a worker that applies the given function to each item in the input channel, without emitting any values.
      log := rivo.Do(func(ctx context.Context, n int) error {
        fmt.Println(n)
        return nil
      })
    
      // `Pipe` composes workers together, returning a new worker
      p := rivo.Pipe3(in, onlyEven, log)
    
      // By passing a context and an input channel to our worker, we can get the output channel.
      // Since our first worker `in` is a generator and does not depend on an input channel, we can pass a nil channel.
      // Also, since log is a sink, we only have to read once from the output channel to know that the pipe has finished.
      <-p(ctx, nil, nil)
    
      // Expected output:
      // 2
      // 4
    }
```

`rivo` provides a set of utilities which can be divided in three main categories:
1. Worker factories: functions that return workers for common use cases, like mapping, filtering, batching, etc.
2. Flow control: functions that help with composing workers together, like `Pipe`, `Merge`, etc.
3. Utilities: functions that help with common tasks, like collecting items from a channel, error handling, etc.


## Worker factories

### Generators
- `Of`: returns a generator that emits the provided values;
- `FromFunc`: returns a generator that emits values returned by the provided function until the function returns false;
- `FromSeq` and `FromSeq2`: return generators that emit the values from the provided iterators;

### Sinks
- `Do`: returns a sink worker that performs a side effect for each item in the input stream;

### Transformers
- `Filter`: returns a worker that filters the input stream using the given function;
- `Map`: returns a worker that applies a function to each item from the input stream;
- `FilterMap`: returns a worker that filters and maps items from the input stream in a single operation;
- `Batch`: returns a worker that groups the input stream into batches of the provided size;
- `Flatten`: returns a worker that flattens the input stream of slices;
- `ForEachOutput`: returns a worker that applies a function to each item, allowing direct output channel access;

Besides these, the library's subdirectories contain more specialized worker factories.

### Package `rivo/io`

- `FromReader`: returns a generator worker that reads from the provided `io.Reader` and emits the read bytes;
- `ToWriter`: returns a sink worker that writes the input stream to the provided `io.Writer`;

### Package `rivo/bufio`

- `FromScanner`: returns a generator worker that reads from the provided `bufio.Scanner` and emits the scanned items;
- `ToWriter`: returns a sink worker that writes the input stream to the provided `bufio.Writer`;

### Package `rivo/csv`

- `FromReader`: returns a generator worker that reads from the provided `csv.Reader` and emits the read records;
- `ToWriter`: returns a sink worker that writes the input stream to the provided `csv.Writer`;

### Configuration Options

Many workers support configuration options to customize their behavior:

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

## Flow control

`rivo` provides functions to compose workers together, allowing you to build complex data processing pipelines:
- `Pipe`, `Pipe2`, `Pipe3`, ... `Pipe10`: compose multiple workers together into a single worker;


## Utilities

`rivo` provides several utility functions to work with streams:

- `Collect`: collects all items from a stream into a slice
- `CollectWithContext`: like `Collect` but respects context cancellation
- `OrDone`: utility function that propagates context cancellation to streams
- `Merge`: merges multiple streams into a single stream

## Examples

More examples can be found in the [examples](./examples) folder.

---

## Contributing

Contributions are welcome! If you have any ideas, suggestions or bug reports, please open an issue or a pull request.

## Roadmap

- [ ] Add more workers, also using the [RxJS list of operators](https://rxjs.dev/guide/operators) as a reference:
  - [x] FilterMap (combines filter and map operations)
  - [x] ForEachOutput (direct output channel access)
  - [ ] Tap (side effects without modifying the stream)
  - [ ] Time-based operators (throttle, debounce, etc.)
  - [ ] SQL-like operators (join, group by, etc.)
- [ ] Add more utilities:
  - [x] Merge (combine multiple streams)
  - [ ] Zip (combine streams element-wise)
  - [ ] Take/Skip operators
- [ ] Performance optimizations and benchmarking
- [ ] Add more examples and tutorials

## License

`rivo` is licensed under the MIT license. See the [LICENSE](./LICENSE) file for details.