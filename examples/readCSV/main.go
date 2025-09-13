package main

import (
	"bufio"
	"context"
	_ "embed"
	"encoding/csv"
	"github.com/agiac/rivo"
	"log"
	"strings"
	"time"

	rivocsv "github.com/agiac/rivo/csv"
)

//go:embed data.csv
var data string

func main() {
	ctx := context.Background()

	r := csv.NewReader(bufio.NewReader(strings.NewReader(data)))
	_, _ = r.Read() // Discard the header line

	readCSV := rivocsv.FromReader(r)

	filterDates := rivo.Filter[[]string](func(ctx context.Context, v []string) bool {
		date, err := time.Parse("2006-01-02", v[5])
		if err != nil {
			return false
		}

		return date.After(time.Date(2023, 1, 20, 0, 0, 0, 0, time.UTC))
	})

	logValues := rivo.Do[[]string](func(ctx context.Context, i []string) {
		log.Println(i)
	})

	<-rivo.Pipe3(
		readCSV,
		filterDates,
		logValues,
	)(ctx, nil, nil)
}
