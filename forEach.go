package rivo

//// ForEach returns a pipeline that applies a function to each item from the input stream.
//// It is intended for side effect and the output stream will only emit the errors returned by the function.
//func ForEach[T any](f func(context.Context, Item[T]) error, opt ...Option) Pipeline[T, struct{}] {
//	forEach := Map[T, struct{}](func(ctx context.Context, item Item[T]) (struct{}, error) {
//		return struct{}{}, f(ctx, item)
//	}, opt...)
//
//	filterNoErrors := Filter(func(ctx context.Context, i Item[struct{}]) (bool, error) {
//		return i.Err != nil, nil
//	})
//
//	return Pipe(forEach, filterNoErrors)
//}
//
//// TODO: review/write tests/remove other ForEach
//
//func ForEach2[T any](f func(context.Context, Item[T]) error, errorHandler Pipeline[struct{}, None], opt ...Option) Pipeline[T, None] {
//	return func(ctx context.Context, in Stream[T]) Stream[None] {
//		out := make(chan Item[None])
//
//		errCh := make(chan Item[struct{}])
//
//		go func() {
//			select {
//			case <-ctx.Done():
//			case <-errorHandler(ctx, errCh):
//			}
//		}()
//
//		go func() {
//			defer close(out)
//			defer close(errCh)
//
//			for item := range OrDone(ctx, in) {
//				err := f(ctx, item)
//				if err != nil {
//					errCh <- Item[struct{}]{Err: err}
//				}
//			}
//		}()
//
//		return out
//	}
//
//}
