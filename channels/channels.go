package channels

import (
	"context"
)

func WriteTo[T any](ctx context.Context, ch chan<- T, value T) (err error) {
	for {
		select {
		case ch <- value:
			goto end
		case <-ctx.Done():
			err = ctx.Err()
			goto end
		}
	}
end:
	return err
}

func ReadFrom[T any](ctx context.Context, ch <-chan T, f func(value T) error) (err error) {
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			goto end
		case value, ok := <-ch:
			if !ok {
				// Channel is closed, we are done
				goto end
			}
			err = f(value)
			if err != nil {
				goto end
			}
		}
	}
end:
	return err
}

func Close[T any](ch chan T) {
	close(ch)
}
