package csv_test

import (
	"context"
	"encoding/csv"
	"github.com/agiac/rivo"
	"strings"
	"testing"

	. "github.com/agiac/rivo/csv"
	"github.com/stretchr/testify/assert"
)

func TestFromReader(t *testing.T) {
	t.Run("read till end of reader", func(t *testing.T) {
		t.Run("without errors", func(t *testing.T) {
			ctx := context.Background()

			r := csv.NewReader(strings.NewReader("1,2,3\n4,5,6\n7,8,9\n"))

			s := FromReader(r)(ctx, nil, nil)

			var rows [][]string
			for item := range s {
				rows = append(rows, item)
			}

			assert.Equal(t, [][]string{
				{"1", "2", "3"},
				{"4", "5", "6"},
				{"7", "8", "9"},
			}, rows)
		})

		t.Run("with errors", func(t *testing.T) {
			ctx := context.Background()

			errs := make(chan error, 1)

			r := csv.NewReader(strings.NewReader("1,2,3\n4,5,6\nerror\n7,8,9\n"))

			s := FromReader(r)(ctx, nil, errs)

			var result [][]string
			for item := range s {
				result = append(result, item)
			}

			close(errs)
			errVals := rivo.Collect(errs)

			assert.Equal(t, [][]string{
				{"1", "2", "3"},
				{"4", "5", "6"},
				{"7", "8", "9"},
			}, result)
			assert.Error(t, errVals[0])
		})

		t.Run("csv reader options", func(t *testing.T) {
			ctx := context.Background()

			r := csv.NewReader(strings.NewReader("1;2;3\n4;5;6\n7;8;9\n"))
			r.Comma = ';'

			s := FromReader(r)(ctx, nil, nil)

			var rows [][]string
			for item := range s {
				rows = append(rows, item)
			}

			assert.Equal(t, [][]string{
				{"1", "2", "3"},
				{"4", "5", "6"},
				{"7", "8", "9"},
			}, rows)
		})

		t.Run("discard header", func(t *testing.T) {
			ctx := context.Background()

			r := csv.NewReader(strings.NewReader("header1,header2,header3\n1,2,3\n4,5,6\n7,8,9\n"))

			s := FromReader(r, DiscardHeader())(ctx, nil, nil)

			var rows [][]string
			for item := range s {
				rows = append(rows, item)
			}

			assert.Equal(t, [][]string{
				{"1", "2", "3"},
				{"4", "5", "6"},
				{"7", "8", "9"},
			}, rows)
		})
	})
}
