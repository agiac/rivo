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

	r := csv.NewReader(bufio.NewReader(strings.NewReader(data)))
	_, _ = r.Read() // Discard the header line

	filter := rivo.FilterFunc[[]string](func(ctx context.Context, i rivo.Item[[]string]) (bool, error) {
		date, err := time.Parse("2006-01-02", i.Val[5])
		if err != nil {
			return false, nil
		}
		return date.After(time.Date(2023, 1, 20, 0, 0, 0, 0, time.UTC)), nil
	})

	log := rivo.DoFunc[[]string](func(ctx context.Context, i rivo.Item[[]string]) {
		log.Println(i.Val)
	})

	<-rivo.Pipe3(rivocsv.FromReader(r), rivo.Filter(filter), rivo.Do(log))(ctx, nil)
}
