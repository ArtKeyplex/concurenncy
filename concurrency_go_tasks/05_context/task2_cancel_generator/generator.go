package generator

import "context"

// Generate возвращает канал, из которого можно читать возрастающие числа,
// начиная с нуля. Генерация прекращается при отмене ctx.
func Generate(ctx context.Context) <-chan int {
	// TODO: реализовать генератор чисел с учётом отмены
	ch := make(chan int)

	go func() {
		defer close(ch)
		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			select {
			case ch <- i:
				i++
			case <-ctx.Done():
				return
			}
		}
	}()

	return ch

}
