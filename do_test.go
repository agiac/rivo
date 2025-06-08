package rivo_test

import (
	"context"
	"errors"
	"fmt"
	. "github.com/agiac/rivo"
)

func ExampleDo() {
	ctx := context.Background()

	in := make(chan Item[int])
	go func() {
		defer close(in)
		in <- Item[int]{Val: 1}
		in <- Item[int]{Val: 2}
		in <- Item[int]{Err: errors.New("error 1")}
		in <- Item[int]{Val: 4}
		in <- Item[int]{Err: errors.New("error 2")}
	}()

	d := Do(func(ctx context.Context, i Item[int]) {
		if i.Err != nil {
			fmt.Printf("ERROR: %v\n", i.Err)
		}
	})

	<-d(ctx, in)

	// Output:
	// ERROR: error 1
	// ERROR: error 2
}
