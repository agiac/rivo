package rivo_test

import (
	"context"
	"fmt"
	"slices"
	"testing"

	. "github.com/agiac/rivo"
	"github.com/stretchr/testify/assert"
)

func ExampleFromSeq() {
	ctx := context.Background()

	seq := slices.Values([]int{1, 2, 3, 4, 5})
	in := FromSeq(seq)

	s := in(ctx, nil)

	for item := range s {
		fmt.Println(item.Val)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
}

func TestFromSeq(t *testing.T) {
	t.Run("create stream from sequence", func(t *testing.T) {
		ctx := context.Background()

		seq := slices.Values([]int{1, 2, 3, 4, 5})
		p := FromSeq(seq)

		got := Collect(p(ctx, nil))

		want := []Item[int]{
			{Val: 1},
			{Val: 2},
			{Val: 3},
			{Val: 4},
			{Val: 5},
		}

		assert.Equal(t, want, got)
	})
}

func ExampleFromSeq2() {
	ctx := context.Background()

	seq := slices.All([]string{"a", "b", "c", "d", "e"})

	in := FromSeq2(seq)

	s := in(ctx, nil)

	for item := range s {
		fmt.Printf("%d, %s\n", item.Val.Val1, item.Val.Val2)
	}

	// Output:
	// 0, a
	// 1, b
	// 2, c
	// 3, d
	// 4, e
}

func TestFromSeq2(t *testing.T) {
	t.Run("create stream from sequence", func(t *testing.T) {
		ctx := context.Background()

		seq := slices.All([]string{"a", "b", "c", "d", "e"})

		p := FromSeq2(seq)

		got := Collect(p(ctx, nil))

		want := []Item[FromSeq2Value[int, string]]{
			{Val: FromSeq2Value[int, string]{Val1: 0, Val2: "a"}},
			{Val: FromSeq2Value[int, string]{Val1: 1, Val2: "b"}},
			{Val: FromSeq2Value[int, string]{Val1: 2, Val2: "c"}},
			{Val: FromSeq2Value[int, string]{Val1: 3, Val2: "d"}},
			{Val: FromSeq2Value[int, string]{Val1: 4, Val2: "e"}},
		}

		assert.Equal(t, want, got)
	})
}
