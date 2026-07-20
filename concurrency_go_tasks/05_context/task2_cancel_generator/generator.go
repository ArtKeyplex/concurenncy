package generator

import (
	"context"
)

// Generate возвращает канал, из которого можно читать возрастающие числа,
// начиная с нуля. Генерация прекращается при отмене ctx.
func Generate(ctx context.Context) <-chan int {
	// TODO: реализовать генератор чисел с учётом отмены
	ch1 := make(chan int)

	select {
	case <-ctx.Done():
		close(ch1)
		return ch1
	default:
	}

	go func() {

		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				close(ch1)
				return

			case ch1 <- i:
			}
		}
	}()

	return ch1
}
