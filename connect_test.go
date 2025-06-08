package rivo_test

import (
	"context"
	"fmt"
	. "github.com/agiac/rivo"
	"github.com/agiac/rivo/core"
	"strings"
)

func ExampleConnect() {
	ctx := context.Background()

	g := Of("Hello", "Hello", "Hello")

	capitalize := Map(func(ctx context.Context, i Item[string]) (string, error) {
		return strings.ToUpper(i.Val), nil
	})

	lowercase := Map(func(ctx context.Context, i Item[string]) (string, error) {
		return strings.ToLower(i.Val), nil
	})

	resA := make([]string, 0)
	a := Do(func(ctx context.Context, i Item[string]) {
		resA = append(resA, i.Val)
	})

	resB := make([]string, 0)
	b := Do(func(ctx context.Context, i Item[string]) {
		resB = append(resB, i.Val)
	})

	p1 := core.Pipe(capitalize, a)
	p2 := core.Pipe(lowercase, b)

	<-core.Pipe(g, Connect(p1, p2))(ctx, nil)

	for _, s := range resA {
		fmt.Println(s)
	}

	for _, s := range resB {
		fmt.Println(s)
	}

	// Output:
	// HELLO
	// HELLO
	// HELLO
	// hello
	// hello
	// hello
}
