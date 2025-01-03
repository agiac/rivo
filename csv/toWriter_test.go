package csv_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"testing"

	. "github.com/agiac/rivo"
	. "github.com/agiac/rivo/csv"
	"github.com/stretchr/testify/assert"
)

func TestToWriter(t *testing.T) {
	ctx := context.Background()

	in := Of([]string{"a", "b", "c"}, []string{"d", "e", "f"}, []string{"g", "h", "i"})

	b := bytes.NewBuffer(nil)
	w := csv.NewWriter(b)

	out := Collect(Pipe(in, ToWriter(w))(ctx, nil))

	assert.Equal(t, "a,b,c\nd,e,f\ng,h,i\n", b.String())
	assert.Equal(t, []Item[struct{}](nil), out)
}
