package core_test

//
//import (
//	"context"
//	"fmt"
//	"github.com/agiac/rivo/core"
//	"sync"
//	"testing"
//
//	"github.com/agiac/rivo"
//	"github.com/stretchr/testify/assert"
//)
//
//func ExampleTee() {
//	ctx := context.Background()
//
//	g := rivo.Of("hello", "hello", "hello")
//
//	out1, out2 := rivo.Tee(g)
//
//	wg := sync.WaitGroup{}
//	wg.Add(2)
//
//	go func() {
//		defer wg.Done()
//		for i := range out1(ctx, nil) {
//			fmt.Println(i.Val)
//		}
//	}()
//
//	go func() {
//		defer wg.Done()
//		for i := range out2(ctx, nil) {
//			fmt.Println(i.Val)
//		}
//	}()
//
//	wg.Wait()
//
//	// Output:
//	// hello
//	// hello
//	// hello
//	// hello
//	// hello
//	// hello
//}
//
//func TestTee(t *testing.T) {
//	ctx := context.Background()
//
//	g := rivo.Of("hello", "hello", "hello")
//
//	out1, out2 := rivo.Tee(g)
//
//	var got1, got2 []rivo.Item[string]
//	wg := sync.WaitGroup{}
//	wg.Add(2)
//
//	go func() {
//		defer wg.Done()
//		got1 = core.Collect(out1(ctx, nil))
//	}()
//
//	go func() {
//		defer wg.Done()
//		got2 = core.Collect(out2(ctx, nil))
//	}()
//
//	wg.Wait()
//
//	want := []rivo.Item[string]{
//		{Val: "hello"},
//		{Val: "hello"},
//		{Val: "hello"},
//	}
//
//	assert.Equal(t, want, got1)
//	assert.Equal(t, want, got2)
//}
//
//func TestTeeN(t *testing.T) {
//	ctx := context.Background()
//
//	g := rivo.Of("hello", "hello", "hello")
//
//	const n = 5
//
//	out := rivo.TeeN(g, n)
//
//	got := make([][]rivo.Item[string], n)
//	wg := sync.WaitGroup{}
//	wg.Add(n)
//	for i := 0; i < n; i++ {
//		go func(i int) {
//			defer wg.Done()
//			got[i] = core.Collect(out[i](ctx, nil))
//		}(i)
//	}
//
//	wg.Wait()
//
//	want := []rivo.Item[string]{
//		{Val: "hello"},
//		{Val: "hello"},
//		{Val: "hello"},
//	}
//
//	for i := 0; i < n; i++ {
//		assert.Equal(t, want, got[i])
//	}
//}
