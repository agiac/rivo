# rivo

`rivo` is a library for stream processing in Go.

## About

`rivo` had two major inspirations:
1. The book ["Concurrency in Go"](https://www.amazon.com/Concurrency-Go-Tools-Techniques-Developers/dp/1491941197);
2. [ReactiveX](https://reactivex.io/), in particular the [Go](https://github.com/ReactiveX/RxGo) and [JS](https://github.com/ReactiveX/rxjs) libraries;

Compared to these sources, `rivo` aims on one side to provide better type safety (both "Concurrency in Go" and RxGo were written in a pre-generics era) and
on the other a slightly more intuitive interface and developer experience (I find ReactiveX hard to adopt, even if rewarding afterward).

## Getting started

### Installation

```shell
  go get github.com/agiac/rivo
```

### Basic concepts

`rivo` has 3 main building blocks: **items**, **streams** and **pipepeables**.

An `Item` is a basic data struct, which contains a value and an optional error. Just like errors are returned next to the result
of a function in synchronous code, so they should be passed along into asynchronous one and handled where more fit.

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

A `Pipeable` is a function that takes a context and a stream of a type and returns a stream of same or different type. It allows to easily compose data pipelines via pipe functions.

```go
type Pipeable[T, U any] func(ctx context.Context, stream Stream[T]) Stream[U]
```

Here's a basic example:

```go
package main

import (
	"fmt"
	"github.com/agiac/rivo"
)

func main() {
	ctx := context.Background()

	// `Of` is a fimple factory function which returs a pipeable which returns a stream that will emit the provided values
	in := rivo.Of(1, 2, 3, 4, 5)

	// `Filter` returns a Pipeable that filters the input stream using the given function.
	onlyEven := rivo.Filter(func(ctx context.Context, i Item[int]) (bool, error) { 
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

	// We can consume the result stream
	for item := range s {
		if item.Err != nil {
			fmt.Printf("ERROR: %v\n", item.Err)
			continue
		}
		fmt.Println(item.Val)
	}
}
```










