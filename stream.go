package rivo

import (
	"fmt"
)

var ErrEOS = fmt.Errorf("end of stream")

type Item[T any] struct {
	Val T
	Err error
}

type Stream[T any] <-chan Item[T]
