package generator

import "context"

// Generate возвращает канал, из которого можно читать возрастающие числа,
// начиная с нуля. Генерация прекращается при отмене ctx.
func Generate(ctx context.Context) <-chan int {
	resChan := make(chan int)

	go func() {
		defer close(resChan)

		n := 0
		for {
			select {
			case <-ctx.Done():
				return
			case resChan <- n:
				n++
			}
		}
	}()

	return resChan
}
