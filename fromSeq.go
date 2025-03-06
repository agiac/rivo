package rivo

import (
	"context"
	"iter"
)

func FromSeq[T any](seq iter.Seq[T]) Pipeline[None, T] {
	return func(ctx context.Context, in Stream[None]) Stream[T] {
		out := make(chan Item[T])

		go func() {
			defer close(out)

			for v := range seq {
				select {
				case <-ctx.Done():
					out <- Item[T]{Err: ctx.Err()}
					return
				case out <- Item[T]{Val: v}:
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

func FromSeq2[T, U any](seq iter.Seq2[T, U]) Pipeline[None, FromSeq2Value[T, U]] {
	return func(ctx context.Context, in Stream[None]) Stream[FromSeq2Value[T, U]] {
		out := make(chan Item[FromSeq2Value[T, U]])

		go func() {
			defer close(out)

			for v1, v2 := range seq {
				select {
				case <-ctx.Done():
					out <- Item[FromSeq2Value[T, U]]{Err: ctx.Err()}
					return
				case out <- Item[FromSeq2Value[T, U]]{Val: FromSeq2Value[T, U]{Val1: v1, Val2: v2}}:
				}
			}
		}()

		return out
	}
}
