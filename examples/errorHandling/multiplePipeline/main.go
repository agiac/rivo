package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"

	rivo "github.com/agiac/rivo/core"
	rivocsv "github.com/agiac/rivo/csv"
)

// This example demonstrates how multiple error handling pipelines can be implemented to run concurrently.

// TODO: consider refactor and better abstractions for error handling pipelines.

func main() {
	ctx := context.Background()

	f, _ := os.Create("examples/errorHandling/multiplePipeline/errors.csv")
	defer f.Close()

	basePipeline := ParseAndDouble()

	valS, errS := rivo.SegregateStream(ctx, basePipeline(ctx, nil), func(i rivo.Item[int]) bool {
		return i.Err == nil
	})

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		<-LogValues()(ctx, valS)
	}()

	go func() {
		defer wg.Done()
		<-rivo.Connect[rivo.Item[int]](
			SaveErrors(f),
			LogErrors(),
		)(ctx, errS)
	}()

	wg.Wait()

	// Expected output (the order might be different because the handleErrors and handleValues pipeline run concurrently):
	// Value: 2
	// Value: 4
	// Error: strconv.Atoi: parsing "3_": invalid syntax
	// Value: 8
	// Value: 10
	// Error: strconv.Atoi: parsing "6**": invalid syntax
	// Error: strconv.Atoi: parsing "?": invalid syntax
	// Value: 16
	// Value: 18
	// Value: 20
}

func ParseAndDouble() rivo.Pipeline[rivo.None, rivo.Item[int]] {
	g := rivo.Of("1", "2", "3_", "4", "5", "6**", "?", "8", "9", "10")

	toInt := rivo.Map(func(ctx context.Context, i string) rivo.Item[int] {
		n, err := strconv.Atoi(i)
		return rivo.Item[int]{Val: n, Err: err}
	})

	double := rivo.Map(func(ctx context.Context, i rivo.Item[int]) rivo.Item[int] {
		if i.Err != nil {
			return i
		}

		return rivo.Item[int]{Val: i.Val * 2}
	})

	return rivo.Pipe3(rivo.Pipeline[rivo.None, string](g), toInt, double)
}

func SaveErrors(f *os.File) rivo.Sync[rivo.Item[int]] {
	toCSVError := rivo.Map[rivo.Item[int], []string](func(ctx context.Context, i rivo.Item[int]) []string {
		return []string{i.Err.Error()}
	})

	save := rivocsv.ToWriter(csv.NewWriter(f))

	logErrors := rivo.Do(func(ctx context.Context, err error) {
		fmt.Printf("Error saving to CSV: %v\n", err)
	})

	return rivo.Sync[rivo.Item[int]](rivo.Pipe3(toCSVError, save, rivo.Pipeline[error, rivo.None](logErrors)))
}

func LogErrors() rivo.Sync[rivo.Item[int]] {
	return rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
		fmt.Printf("Error: %v\n", i.Err)
	})
}

func LogValues() rivo.Sync[rivo.Item[int]] {
	return rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
		fmt.Printf("Value: %d\n", i.Val)
	})
}
