package main

import (
	"bufio"
	"context"
	_ "embed"
	"encoding/csv"
	"log"
	"strings"
	"time"

	"github.com/agiac/rivo"
	rivocsv "github.com/agiac/rivo/csv"
)

//go:embed data.csv
var data string

func main() {
	ctx := context.Background()

	// Create a new CSV reader and discard the first line
	r := csv.NewReader(bufio.NewReader(strings.NewReader(data)))
	_, _ = r.Read()

	p := rivo.Pipe(rivocsv.FromReader(r), rivo.Filter(func(ctx context.Context, i rivo.Item[[]string]) (bool, error) {
		date, err := time.Parse("2006-01-02", i.Val[5])
		if err != nil {
			return false, nil
		}
		return date.After(time.Date(2023, 1, 20, 0, 0, 0, 0, time.UTC)), nil
	}))

	for item := range p(ctx, nil) {
		if item.Err != nil {
			log.Printf("ERROR: %v\n", item.Err)
			continue
		}
		log.Println(item.Val)
	}
}
