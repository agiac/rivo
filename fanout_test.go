package rivo

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func ExampleFanout() {
	ctx := context.Background()

	g := Of("Hello", "Hello", "Hello")

	capitalize := Map(func(ctx context.Context, i string) (string, error) {
		return strings.ToUpper(i), nil
	})

	lowercase := Map(func(ctx context.Context, i string) (string, error) {
		return strings.ToLower(i), nil
	})

	resA := make([]string, 0)
	a := Do(func(ctx context.Context, i string) error {
		resA = append(resA, i)
		return nil
	})

	resB := make([]string, 0)
	b := Do(func(ctx context.Context, i string) error {
		resB = append(resB, i)
		return nil
	})

	p1 := Pipe(capitalize, a)
	p2 := Pipe(lowercase, b)

	<-Fanout(g, p1, p2)(ctx, nil, nil)

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

func TestFanout(t *testing.T) {
	t.Skip("TODO") // TODO
}
