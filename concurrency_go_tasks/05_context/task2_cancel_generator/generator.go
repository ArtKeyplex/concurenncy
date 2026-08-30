package generator

import "context"

// Generate возвращает канал, из которого можно читать возрастающие числа,
// начиная с нуля. Генерация прекращается при отмене ctx.
func Generate(ctx context.Context) <-chan int {

	out := make(chan int)

	go func() {
		defer close(out)
		n := 0
		for {
			if ctx.Err() != nil {
				return
			}

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
