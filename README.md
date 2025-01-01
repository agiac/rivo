# rivo

`rivo` is a library for stream processing in Go.

## About

`rivo` had two major inspirations:
1. The book ["Concurrency in Go"](https://www.amazon.com/Concurrency-Go-Tools-Techniques-Developers/dp/1491941197);
2. [ReactiveX](https://reactivex.io/), in particular the [Go](https://github.com/ReactiveX/RxGo) and [JS](https://github.com/ReactiveX/rxjs) libraries;

Compared to these sources, `rivo` aims to provide better type safety (both "Concurrency in Go" and RxGo were written in a pre-generics era and make a heavy use of `interface{}`) 
and a more intuitive API and developer experience.

## Getting started

### Installation

```shell
  go get github.com/agiac/rivo
```

### Basic concepts

`rivo` has 3 main building blocks: **items**, **streams** and **pipepeables**.

An `Item` is a basic struct, which contains a value and an optional error. Just like errors are returned next to the result 
of a function in synchronous code, they should be passed along into asynchronous code and handled where more appropriate.

```go
type Item[T any] struct {
	Val T
	Err error
}
```

A `Stream` is a read only channel of items.

```go
type Stream[T any] <-chan Item[T]
```

A `Pipeable` is a function that takes a context and a stream of one type and returns a stream of the same or a different type. 
It allows for easy composition of data pipelines via pipe functions.

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

Many pipeables accepts a common set of optional parameters. These can be provided via functional options.

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/agiac/rivo"
)

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
```

The currently available options are:

- `WithPoolSize(int)`: sets the number of goroutines that will be used to process items. Default is 1.
- `WithBufferSize(int)`: sets the buffer size of the output channel. Default is 0;
- `WithStopOnError(bool)`: if true, the pipeable will stop processing items when an error is encountered. Default is false.

More examples can be found in the [examples](./examples) folder.

## Roadmap

- [ ] Add more pipeables, also using the [RxJS list of operators](https://rxjs.dev/guide/operators) as a reference:
  - [ ] IO
  - [ ] Error handling
  - [ ] Time-based
  - [ ] SQL
  - [ ] AWS
  - [ ] ...
- [ ] Add more examples

## License

`rivo` is licensed under the MIT license. See the [LICENSE](./LICENSE) file for details.









