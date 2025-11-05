package generator

import "context"

// Generate возвращает канал, из которого можно читать возрастающие числа,
// начиная с нуля. Генерация прекращается при отмене ctx.
func Generate(ctx context.Context) <-chan int {
	// TODO: реализовать генератор чисел с учётом отмены
	ch := make(chan int)
	select {
	case <-ctx.Done():
		close(ch)
		return ch
	default:

	}
	go func() {
		num := 0
		defer close(ch)
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
