package generator

import "context"

// Generate возвращает канал, из которого можно читать возрастающие числа,
// начиная с нуля. Генерация прекращается при отмене ctx.
func Generate(ctx context.Context) <-chan int {
	// TODO: реализовать генератор чисел с учётом отмены

	ch := make(chan int)
	if ctx.Err() != nil {
		close(ch)
		return (ch)
	}

	go func() {
		defer close(ch)
		num := 0

		for {
			select {
			case <-ctx.Done():
				return
			case ch <- num:
				num++
			}
		}
	}()
	return ch

}
