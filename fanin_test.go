package rivo

import (
	"context"
	"fmt"
	"slices"
	"testing"
)

func ExampleFanin() {
	ctx := context.Background()

	g1 := Of("Hello", "World")
	g2 := Of("Foo", "Bar")

	gg := Fanin(g1, g2)

	res := Collect(gg(ctx, nil, nil))

	slices.Sort(res)

	fmt.Println(res)

	// Output:
	// [Bar Foo Hello World]
}

func TestFanin(t *testing.T) {
	t.Skip("TODO") // TODO
}
