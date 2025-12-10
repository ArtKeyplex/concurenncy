package generator

import "context"

// Generate возвращает канал, из которого можно читать возрастающие числа,
// начиная с нуля. Генерация прекращается при отмене ctx.
func Generate(ctx context.Context) <-chan int {
	// TODO: реализовать генератор чисел с учётом отмены
	out := make(chan int)

	go func() {
		defer close(out)
		for i := 0; ; i++ {
			if ctx.Err() != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case out <- i:

			}
		}
	}()
	/*/go func() {
		<-ctx.Done()
		close(out)
	}() то проходит тесты, то нет /*/
	return out
}
