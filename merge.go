package rivo

import (
	"context"
	"sync"
)

// Merge merges multiple input channels into a single output channel.
// It stops merging when the context is cancelled or all input channels are closed.
func Merge[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		wg := sync.WaitGroup{}

		for _, ch := range channels {
			wg.Add(1)
			go func(c <-chan T) {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					case item, ok := <-c:
						if !ok {
							return
						}
						select {
						case <-ctx.Done():
							return
						case out <- item:
						}
					}
				}
			}(ch)
		}

		wg.Wait()

	}()

	return out
}
