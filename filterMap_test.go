package rivo_test

import (
	"context"
	"fmt"
	. "github.com/agiac/rivo"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleFilterMap() {
	ctx := context.Background()

	in := Of(1, 2, 3, 4, 5)

	// Filter even numbers and multiply by 10
	filterMapEvenAndMultiply := FilterMap(func(ctx context.Context, n int) (bool, int) {
		if n%2 == 0 {
			return true, n * 10
		}
		return false, 0
	})

	p := Pipe(in, filterMapEvenAndMultiply)

	s := p(ctx, nil)

	for item := range s {
		fmt.Println(item)
	}

	// Output:
	// 20
	// 40
}

func TestFilterMap(t *testing.T) {
	filterMapFunc := func(ctx context.Context, i int) (bool, string) {
		if i%2 == 0 {
			return true, fmt.Sprintf("even-%d", i)
		}
		return false, ""
	}

	t.Run("filter and map items", func(t *testing.T) {
		ctx := context.Background()

		g := Of(1, 2, 3, 4, 5)
		fm := FilterMap(filterMapFunc)

		got := Collect(Pipe(g, fm)(ctx, nil))
		want := []string{"even-2", "even-4"}

		assert.Equal(t, want, got)
	})

	t.Run("with context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		g := Of(1, 2, 3, 4, 5)
		fm := FilterMap(filterMapFunc)

		got := Collect(fm(ctx, g(ctx, nil)))

		assert.Lessf(t, len(got), 5, "expected less than 5 items, got %d", len(got))
	})

	t.Run("with buffer size", func(t *testing.T) {
		ctx := context.Background()

		g := Of(1, 2, 3, 4, 5)
		fm := FilterMap(filterMapFunc, FilterMapBufferSize(3))

		got := Collect(fm(ctx, g(ctx, nil)))
		want := []string{"even-2", "even-4"}

		assert.Equal(t, want, got)
	})

	t.Run("with pool size", func(t *testing.T) {
		ctx := context.Background()

		g := Of(1, 2, 3, 4, 5)
		fm := FilterMap(filterMapFunc, FilterMapPoolSize(3))

		got := Collect(Pipe(g, fm)(ctx, nil))
		want := []string{"even-2", "even-4"}

		assert.ElementsMatch(t, want, got) // Order might not be guaranteed with pool size > 1
	})
}

func TestFilterMapValues(t *testing.T) {
	t.Run("filters out errors", func(t *testing.T) {
		ctx := context.Background()
		in := Of(
			Item[int]{Val: 1},
			Item[int]{Err: fmt.Errorf("error 1")},
			Item[int]{Val: 3},
			Item[int]{Err: fmt.Errorf("error 2")},
			Item[int]{Val: 5},
		)

		fmv := FilterMapValues[int]()
		p := Pipe(in, fmv)
		got := Collect(p(ctx, nil))
		want := []int{1, 3, 5}

		assert.Equal(t, want, got)
	})

	t.Run("empty input", func(t *testing.T) {
		ctx := context.Background()
		in := Of[Item[int]]()

		fmv := FilterMapValues[int]()
		p := Pipe(in, fmv)
		got := Collect(p(ctx, nil))
		want := []int(nil)

		assert.Equal(t, want, got)
	})

	t.Run("all errors", func(t *testing.T) {
		ctx := context.Background()
		in := Of(
			Item[int]{Err: fmt.Errorf("error 1")},
			Item[int]{Err: fmt.Errorf("error 2")},
		)

		fmv := FilterMapValues[int]()
		p := Pipe(in, fmv)
		got := Collect(p(ctx, nil))
		want := []int(nil)

		assert.Equal(t, want, got)
	})

	t.Run("all values", func(t *testing.T) {
		ctx := context.Background()
		in := Of(
			Item[int]{Val: 1},
			Item[int]{Val: 2},
		)

		fmv := FilterMapValues[int]()
		p := Pipe(in, fmv)
		got := Collect(p(ctx, nil))
		want := []int{1, 2}

		assert.Equal(t, want, got)
	})
}

func TestFilterMapErrors(t *testing.T) {
	err1 := fmt.Errorf("error 1")
	err2 := fmt.Errorf("error 2")

	t.Run("filters out values", func(t *testing.T) {
		ctx := context.Background()
		in := Of(
			Item[int]{Val: 1},
			Item[int]{Err: err1},
			Item[int]{Val: 3},
			Item[int]{Err: err2},
			Item[int]{Val: 5},
		)

		fme := FilterMapErrors[int]()
		p := Pipe(in, fme)
		got := Collect(p(ctx, nil))
		want := []error{err1, err2}

		assert.Equal(t, want, got)
	})

	t.Run("empty input", func(t *testing.T) {
		ctx := context.Background()
		in := Of[Item[int]]()

		fme := FilterMapErrors[int]()
		p := Pipe(in, fme)
		got := Collect(p(ctx, nil))
		want := []error(nil)

		assert.Equal(t, want, got)
	})

	t.Run("all values", func(t *testing.T) {
		ctx := context.Background()
		in := Of(
			Item[int]{Val: 1},
			Item[int]{Val: 2},
		)

		fme := FilterMapErrors[int]()
		p := Pipe(in, fme)
		got := Collect(p(ctx, nil))
		want := []error(nil)

		assert.Equal(t, want, got)
	})

	t.Run("all errors", func(t *testing.T) {
		ctx := context.Background()
		in := Of(
			Item[int]{Err: err1},
			Item[int]{Err: err2},
		)

		fme := FilterMapErrors[int]()
		p := Pipe(in, fme)
		got := Collect(p(ctx, nil))
		want := []error{err1, err2}

		assert.Equal(t, want, got)
	})
}
