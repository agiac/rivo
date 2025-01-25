package rivo_test

//
//import (
//	"context"
//	"fmt"
//	"strconv"
//	"testing"
//
//	. "github.com/agiac/rivo"
//	"github.com/stretchr/testify/assert"
//)
//
//func ExampleForEach() {
//	ctx := context.Background()
//
//	g := Of(1, 2, 3, 4, 5)
//
//	f := ForEach(func(ctx context.Context, i Item[int]) error {
//		// Do some side effect
//		// ...
//		// Simulate an error
//		if i.Val == 3 {
//			return fmt.Errorf("an error")
//		}
//
//		return nil
//	})
//
//	s := Pipe(g, f)(ctx, nil)
//
//	for item := range s {
//		fmt.Printf("item: %v; error: %v\n", item.Val, item.Err)
//	}
//
//	// Output:
//	// item: {}; error: an error
//}
//
//func ExampleForEach2() {
//	ctx := context.Background()
//
//	g := Of("1", "!2", "w", "4", "5")
//
//	errHandler := Do[struct{}](func(ctx context.Context, i Item[struct{}]) {
//		fmt.Printf("Error\n")
//	})
//
//	res := make([]int, 0)
//	forEachFn := func(ctx context.Context, i Item[string]) error {
//		n, err := strconv.Atoi(i.Val)
//		if err != nil {
//			return err
//		}
//
//		res = append(res, n)
//
//		return nil
//	}
//
//	f := ForEach2[string](forEachFn, errHandler)
//
//	<-Pipe(g, f)(ctx, nil)
//
//	for _, n := range res {
//		fmt.Printf("%v\n", n)
//	}
//
//	// Output:
//	// Error
//	// Error
//	// 1
//	// 4
//	// 5
//}
//
//func TestForEach(t *testing.T) {
//	t.Run("for each item", func(t *testing.T) {
//		ctx := context.Background()
//
//		g := Of(1, 2, 3, 4, 5)
//
//		sideEffect := make([]int, 0)
//		f := ForEach(func(ctx context.Context, i Item[int]) error {
//			sideEffect = append(sideEffect, i.Val)
//			return nil
//		})
//
//		errs := Collect(Pipe(g, f)(ctx, nil))
//
//		assert.Equal(t, []int{1, 2, 3, 4, 5}, sideEffect)
//		assert.Equal(t, 0, len(errs))
//	})
//
//	t.Run("forward errors", func(t *testing.T) {
//		ctx := context.Background()
//
//		g := Of(1, 2, 3, 4, 5)
//
//		sideEffect := make([]int, 0)
//		f := ForEach(func(ctx context.Context, i Item[int]) error {
//			sideEffect = append(sideEffect, i.Val)
//
//			if i.Val == 3 {
//				return fmt.Errorf("an error")
//			}
//
//			return nil
//		})
//
//		errs := Collect(Pipe(g, f)(ctx, nil))
//
//		assert.Equal(t, []int{1, 2, 3, 4, 5}, sideEffect)
//		assert.Equal(t, 1, len(errs))
//		assert.Error(t, errs[0].Err)
//	})
//}
