package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/agiac/rivo"
	rivocsv "github.com/agiac/rivo/csv"
)

// This example demonstrates how multiple error handling pipelines can be implemented to run concurrently.

func main() {
	ctx := context.Background()

	f, _ := os.Create("examples/errorHandling/multiplePipeline/errors.csv")
	defer f.Close()

	vals, errs := rivo.SegregateErrors(ParseAndDouble())(ctx, nil)

	errs1, errs2 := rivo.Tee(errs)(ctx, nil)

	<-rivo.Connect(
		rivo.Pipe(vals, LogValues()),
		rivo.Pipe(errs1, SaveErrors(f)),
		rivo.Pipe(errs2, LogErrors()),
	)(ctx, nil)

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

func ParseAndDouble() rivo.Pipeline[rivo.None, int] {
	g := rivo.Of("1", "2", "3_", "4", "5", "6**", "?", "8", "9", "10")

	toInt := rivo.Map(func(ctx context.Context, i rivo.Item[string]) (int, error) {
		if i.Err != nil {
			return 0, i.Err // Pass errors along
		}

		return strconv.Atoi(i.Val)
	})

	double := rivo.Map(func(ctx context.Context, i rivo.Item[int]) (int, error) {
		if i.Err != nil {
			return 0, i.Err // Pass errors along
		}

		return i.Val * 2, nil
	})

	return rivo.Pipe3(g, toInt, double)
}

func SaveErrors(f *os.File) rivo.Pipeline[int, rivo.None] {
	toCSVError := rivo.Map[int, []string](func(ctx context.Context, i rivo.Item[int]) ([]string, error) {
		return []string{i.Err.Error()}, nil
	})

	save := rivocsv.ToWriter(csv.NewWriter(f))

	logErrors := rivo.Do(func(ctx context.Context, i rivo.Item[struct{}]) {
		fmt.Printf("Error saving to CSV: %v\n", i.Err)
	})

	return rivo.Pipe3(toCSVError, save, logErrors)
}

func LogErrors() rivo.Pipeline[int, rivo.None] {
	return rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
		fmt.Printf("Error: %v\n", i.Err)
	})
}

func LogValues() rivo.Pipeline[int, rivo.None] {
	return rivo.Do(func(ctx context.Context, i rivo.Item[int]) {
		fmt.Printf("Value: %d\n", i.Val)
	})
}
