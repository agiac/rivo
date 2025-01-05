package rivo

import (
	"context"
	"iter"
)

func FromSeq[T any](seq iter.Seq[T], opt ...Option) Generator[T] {
	next, stop := iter.Pull(seq)
	return FromFunc[T](
		func(ctx context.Context) (T, error) {
			v, ok := next()
			if !ok {
				return v, ErrEOS
			}

			return v, nil
		},
		append(opt, WithOnBeforeClose(func(ctx context.Context) error {
			stop()
			return nil
		}))...,
	)
}

type FromSeq2Value[T, U any] struct {
	Val1 T
	Val2 U
}

func FromSeq2[T, U any](seq iter.Seq2[T, U], opts ...Option) Generator[FromSeq2Value[T, U]] {
	next, stop := iter.Pull2(seq)
	return FromFunc[FromSeq2Value[T, U]](
		func(ctx context.Context) (FromSeq2Value[T, U], error) {
			v1, v2, ok := next()
			if !ok {
				return FromSeq2Value[T, U]{Val1: v1, Val2: v2}, ErrEOS
			}

			return FromSeq2Value[T, U]{Val1: v1, Val2: v2}, nil
		},
		append(opts, WithOnBeforeClose(func(ctx context.Context) error {
			stop()
			return nil
		}))...,
	)
}
