package rivo_test

import (
	"context"
	"fmt"
	"testing"

	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

func ExampleFork() {
	ctx := context.Background()

	g := Of(1, 2, 3)

	g1, g2 := Fork(g)

	for i := range g1(ctx, nil) {
		fmt.Printf("g1: %v\n", i.Val)
	}

	for i := range g2(ctx, nil) {
		fmt.Printf("g2: %v\n", i.Val)
	}

	// Output:
	// g1: 1
	// g1: 2
	// g1: 3
	// g2: 1
	// g2: 2
	// g2: 3
}

func TestForkN(t *testing.T) {
	ctx := context.Background()

	const n = 3

	g := Of(1, 2, 3)

	gs := ForkN(g, n)

	res := make([][]Item[int], n)
	for i := range n {
		res[i] = Collect((gs[i])(ctx, nil))
	}

	expected := [][]Item[int]{
		{{Val: 1}, {Val: 2}, {Val: 3}},
		{{Val: 1}, {Val: 2}, {Val: 3}},
		{{Val: 1}, {Val: 2}, {Val: 3}},
	}

	for i := range n {
		assert.Equal(t, expected[i], res[i])
	}
}
