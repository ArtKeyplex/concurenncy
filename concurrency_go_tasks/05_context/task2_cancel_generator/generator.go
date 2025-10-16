package generator

import "context"

// Generate возвращает канал, из которого можно читать возрастающие числа,
// начиная с нуля. Генерация прекращается при отмене ctx.
func Generate(ctx context.Context) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		number := 0

		for {
			select {
			case <-ctx.Done():
				return
			case ch <- number:
				number++
			}
		}
	}()

	return ch
}
