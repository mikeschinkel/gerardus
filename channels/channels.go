package channels

import (
	"context"
)

// WriteTo writes an item to a channel of that items type while allowing for cancellation.
// e.g.
//
//			// See https://goplay.tools/snippet/bBq3PpAeMa0
//			intChan := chan int
//			ctx,cancel := context.WithCancel(context.Background())
//			err := channels.WriteTo(ctx,intChan,10)
//		  if err != nil {
//	       cancel()
//					return err
//			}
//	   close(intChan)
//
// In the case of a channel of an interface, e
//
//	//See https://goplay.tools/snippet/m6M8-GwfyUq
//	type Suiter interface {
//		Suit()
//	}
//	type Hearts struct{}
//	func (Hearts) Suit() {}
//	func Example() {
//		suitChan := make(chan Suiter, 10)
//		ctx := context.Background()
//		err = WriteTo[Suiter](ctx, suitChan, Hearts{})
//		if err != nil {
//				cancel()
//				return err
//		}
//		close(suitChan)
//	}
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

// ReadAllFrom reads all items from a channel and calls the closure for each item.
// e.g.
//
//	intChan := make(chan int, 10)
//	// Do something to fill intChan
//	ctx, cancel := context.WithCancel(context.Background())
//	err := ReadFrom(ctx, intChan, func(i int) error {
//		if i%2 == 0 { // Simulate an error
//			return fmt.Errorf("Oops! We can't do even: %d", i)
//			cancel()
//		}
//		fmt.Println(i)
//		return nil
//	})
//	if err != nil {
//			cancel()
//			return err
//	}
func ReadAllFrom[T any](ctx context.Context, ch <-chan T, f func(value T) error) (err error) {
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

// CancelOnErr allows calling a closure that may return an error, and if it does
// call the cancel function passed in. Useful when used repeatedly in the same func.
func CancelOnErr(cancel context.CancelFunc, f func() error) func() error {
	return func() error {
		err := f()
		if err != nil {
			cancel()
		}
		return err
	}
}
