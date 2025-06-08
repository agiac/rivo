package rivo_test

//import (
//	"context"
//	"fmt"
//	"github.com/agiac/rivo/core"
//	"strconv"
//	"sync"
//	"testing"
//
//	"github.com/agiac/rivo"
//	"github.com/stretchr/testify/assert"
//)
//
//func ExampleSegregate() {
//	ctx := context.Background()
//
//	g := rivo.Of("1", "2", "3", "4", "5")
//
//	toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
//		return strconv.Atoi(i.Val)
//	})
//
//	p := core.Pipe(g, toInt)
//
//	even, odd := rivo.Segregate(p, func(item rivo.Item[int]) bool {
//		return item.Val%2 == 0
//	})
//
//	evens := make([]int, 0)
//	odds := make([]int, 0)
//
//	<-rivo.Connect(
//		core.Pipe(even, rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
//			evens = append(evens, i.Val)
//		})),
//		core.Pipe(odd, rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
//			odds = append(odds, i.Val)
//		})),
//	)(ctx, nil)
//
//	for _, i := range append(evens, odds...) {
//		fmt.Println(i)
//	}
//
//	// Output:
//	// 2
//	// 4
//	// 1
//	// 3
//	// 5
//}
//
//func TestSegregate(t *testing.T) {
//	t.Run("all values", func(t *testing.T) {
//
//		ctx := context.Background()
//
//		g := rivo.Of("1", "2", "3", "4", "5", "6", "7", "8", "9", "10")
//
//		toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
//			return strconv.Atoi(i.Val)
//		})
//
//		p := core.Pipe(g, toInt)
//
//		even, odd := rivo.Segregate(p, func(item rivo.Item[int]) bool {
//			return item.Val%2 == 0
//		})
//
//		evens := make([]int, 0)
//		odds := make([]int, 0)
//
//		<-rivo.Connect(
//			core.Pipe(even, rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
//				evens = append(evens, i.Val)
//			})),
//			core.Pipe(odd, rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
//				odds = append(odds, i.Val)
//			})),
//		)(ctx, nil)
//
//		assert.Equal(t, []int{2, 4, 6, 8, 10}, evens)
//		assert.Equal(t, []int{1, 3, 5, 7, 9}, odds)
//	})
//
//	t.Run("context cancellation", func(t *testing.T) {
//		ctx, cancel := context.WithCancel(context.Background())
//		cancel()
//
//		g := rivo.Of("1", "2", "3", "4", "5", "6", "7", "8", "9", "10")
//
//		toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
//			return strconv.Atoi(i.Val)
//		})
//
//		p := core.Pipe(g, toInt)
//
//		even, odd := rivo.Segregate(p, func(item rivo.Item[int]) bool {
//			return item.Val%2 == 0
//		})
//
//		evens := make([]int, 0)
//		odds := make([]int, 0)
//
//		wg := sync.WaitGroup{}
//		wg.Add(2)
//		go func() {
//			defer wg.Done()
//			for i := range even(ctx, nil) {
//				evens = append(evens, i.Val)
//			}
//		}()
//
//		go func() {
//			defer wg.Done()
//			for i := range odd(ctx, nil) {
//				odds = append(odds, i.Val)
//			}
//		}()
//
//		wg.Wait()
//
//		assert.Empty(t, evens)
//		assert.Empty(t, odds)
//	})
//}
