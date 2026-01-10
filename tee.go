package rivo

import (
	"context"
)

// TODO: tests

// Tee returns 2 channels that each receive a copy of each item from the input channel.
func Tee[T any](ctx context.Context, in <-chan T) (<-chan T, <-chan T) {
	ss := TeeN(ctx, in, 2)
	return ss[0], ss[1]
}

// TeeN returns n channels that each receive a copy of each item from the input channel.
func TeeN[T any](ctx context.Context, in <-chan T, n int) []<-chan T {
	if n <= 0 {
		panic("n must be greater than 0")
	}

	out := make([]chan T, n)
	for i := 0; i < n; i++ {
		out[i] = make(chan T)
	}

	go func() {
		defer func() {
			for i := 0; i < n; i++ {
				close(out[i])
			}
		}()

		for item := range OrDone(ctx, in) {
			for i := 0; i < n; i++ {
				select {
				case <-ctx.Done():
					return
				case out[i] <- item:
				}
			}
		}
	}()

	ins := make([]<-chan T, n)
	for i := 0; i < n; i++ {
		ins[i] = out[i]
	}

	return ins
}
