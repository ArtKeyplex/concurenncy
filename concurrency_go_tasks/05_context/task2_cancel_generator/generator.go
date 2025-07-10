package generator

import "context"

func Generate(ctx context.Context) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		n := 0
		for {
			select {
			case <-ctx.Done():
				return
			case out <- n:
				n++
			}
		}
	}()

	return out
}