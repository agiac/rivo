package csv_test

import (
	"context"
	"encoding/csv"
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

			s := FromReader(r)(ctx, nil)

			var rows [][]string
			for item := range s {
				assert.NoError(t, item.Err)
				rows = append(rows, item.Val)
			}

			assert.Equal(t, [][]string{
				{"1", "2", "3"},
				{"4", "5", "6"},
				{"7", "8", "9"},
			}, rows)
		})

		t.Run("with errors", func(t *testing.T) {
			ctx := context.Background()

			r := csv.NewReader(strings.NewReader("1,2,3\n4,5,6\nerror\n7,8,9\n"))

			s := FromReader(r)(ctx, nil)

			var result [][]string
			var errs []error
			for item := range s {
				if item.Err != nil {
					errs = append(errs, item.Err)
				} else {
					result = append(result, item.Val)
				}
			}

			assert.Equal(t, [][]string{
				{"1", "2", "3"},
				{"4", "5", "6"},
				{"7", "8", "9"},
			}, result)
			assert.Error(t, errs[0])
		})

		t.Run("csv reader options", func(t *testing.T) {
			ctx := context.Background()

			r := csv.NewReader(strings.NewReader("1;2;3\n4;5;6\n7;8;9\n"))
			r.Comma = ';'

			s := FromReader(r)(ctx, nil)

			var rows [][]string
			for item := range s {
				if item.Err != nil {
					assert.Fail(t, "unexpected error", item.Err)
					return
				}
				rows = append(rows, item.Val)
			}

			assert.Equal(t, [][]string{
				{"1", "2", "3"},
				{"4", "5", "6"},
				{"7", "8", "9"},
			}, rows)
		})
	})
}
