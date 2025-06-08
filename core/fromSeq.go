package core

import (
	"context"
	"iter"
)

func FromSeq[T any](seq iter.Seq[T]) Generator[T] {
	return func(ctx context.Context, in Stream[None]) Stream[T] {
		out := make(chan T)

		go func() {
			defer close(out)

			for v := range seq {
				select {
				case <-ctx.Done():
					return
				case out <- v:
				}
			}
		}()

		return out
	}
}

type FromSeq2Value[T, U any] struct {
	Val1 T
	Val2 U
}

func FromSeq2[T, U any](seq iter.Seq2[T, U]) Generator[FromSeq2Value[T, U]] {
	return func(ctx context.Context, in Stream[None]) Stream[FromSeq2Value[T, U]] {
		out := make(chan FromSeq2Value[T, U])

		go func() {
			defer close(out)

			for v1, v2 := range seq {
				select {
				case <-ctx.Done():
					return
				case out <- FromSeq2Value[T, U]{Val1: v1, Val2: v2}:
				}
			}
		}()

		return out
	}
}
