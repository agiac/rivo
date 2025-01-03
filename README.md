# rivo

`rivo` is a library for stream processing in Go.

## About

`rivo` had two major inspirations:
1. The book ["Concurrency in Go"](https://www.amazon.com/Concurrency-Go-Tools-Techniques-Developers/dp/1491941197);
2. [ReactiveX](https://reactivex.io/), in particular the [Go](https://github.com/ReactiveX/RxGo) and [JS](https://github.com/ReactiveX/rxjs) libraries;

Compared to these sources, `rivo` aims to provide better type safety (both "Concurrency in Go" and RxGo were written in a pre-generics era and make a heavy use of `interface{}`) 
and a more intuitive API and developer experience (Rx is very powerful, but can be overwhelming for newcomers).

## Getting started

### Installation

```shell
  go get github.com/agiac/rivo
```

### Basic concepts

`rivo` has 3 main types, which are the building blocks of the library: `Item`, `Stream` and `Pipeable`.

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

`Pipeable` is a function that takes a `context.Context` and a `Stream` of one type and returns a `Stream` of the same or a different type.
Pipeables can be composed together using the one of the `Pipe` functions.

```go
type Pipeable[T, U any] func(ctx context.Context, stream Stream[T]) Stream[U]
```

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

	// `Of` is a factory function which returns a pipeable which returns a stream that will emit the provided values
	in := rivo.Of(1, 2, 3, 4, 5)

	// `Filter` returns a pipeable that filters the input stream using the given function.
	onlyEven := rivo.Filter(func(ctx context.Context, i rivo.Item[int]) (bool, error) {
		// Always check for errors
		if i.Err != nil {
			return true, i.Err // Propagate the error
		}

		return i.Val%2 == 0, nil
	})

	// `Pipe` composes pipeables together, returning a new pipeable
	p := rivo.Pipe(in, onlyEven)

	// By passing a context and an input channel to our pipeable, we can get the output stream.
	// Since our first pipeable `in` does not depend on an input stream, we pass a nil channel.
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
```
## Pipeable factories

`rivo` comes with a set of built-in pipeable factories. The core ones are:

- `Of`: returns a pipeable which returns a stream that will emit the provided values;
- `FromFunc`: returns a pipeable which returns a stream that will emit the values returned by the provided function;
- `Filter`: returns a pipeable which filters the input stream using the given function;
- `Map`: returns a pipeable which maps the input stream using the given function;
- `ForEach`: returns a pipeable which applies the given function to each item in the input stream and forwards only the errors;
- `Do`: returns a pipeable which performs a side effect for each item in the input stream;
- `Parallel`: returns a pipeable which applies the given pipeables to the input stream concurrently;

Besides these, the directories of the library contain more specialized pipeables factories.

### Package `rivo/io`

- `FromReader`: returns a pipeable which reads from the provided `io.Reader` and emits the read bytes;
- `ToWriter`: returns a pipeable which writes the input stream to the provided `io.Writer`;

### Package `rivo/bufio`

- `FromScanner`: returns a pipeable which reads from the provided `bufio.Scanner` and emits the scanned items;
- `ToScanner`: returns a pipeable which writes the input stream to the provided `bufio.Writer`;

### Package `rivo/csv`

- `FromReader`: returns a pipeable which reads from the provided `csv.Reader` and emits the read records;
- `ToWriter`: returns a pipeable which writes the input stream to the provided `csv.Writer`;

### Optional parameters

Many pipeable factories accepts a common set of optional parameters. These can be provided via functional options.

```go
  double := rivo.Map(
	  func(ctx context.Context, i rivo.Item[int]) (int, error) { return i.Val * 2, nil  },
	  // `Pass additional options to the pipeable
	  rivo.WithBufferSize(1), 
	  rivo.WithPoolSize(runtime.NumCPU()), 
	  )
```

The currently available options are:

- `WithPoolSize(int)`: sets the number of goroutines that will be used to process items. Default is 1.
- `WithBufferSize(int)`: sets the buffer size of the output channel. Default is 0;
- `WithStopOnError(bool)`: if true, the pipeable will stop processing items when an error is encountered. Default is false.
- `WithOnBeforeClosed(func(context.Context) error)`: a function that will be called before the output channel is closed.

## Utilities

`rivo` also comes with a set of utilities which cannot be expressed as pipeables but can be useful when working with streams:

- `OrDone`: returns a channel which will be closed when the provided context is done;
- `Tee` and `TeeN`: returns n streams that each receive a copy of each item from the input stream;


## Error handling

TODO

## Examples

More examples can be found in the [examples](./examples) folder.

---

## Contributing

Contributions are welcome! If you have any ideas, suggestions or bug reports, please open an issue or a pull request.

## Roadmap

- [ ] Generator and Sync type aliases, once bug in Go 1.23 is fixed
- [ ] Add more pipeables, also using the [RxJS list of operators](https://rxjs.dev/guide/operators) as a reference:
  - [ ] Batch
  - [ ] Error handling
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









