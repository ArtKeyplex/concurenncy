package fibonacci

// Fib возвращает канал, из которого можно читать первые n чисел Фибоначчи.
func Fib(n int) <-chan int {
	ch := make(chan int)

	go func() {
		x, y := 0, 1

		for i := 0; i < n; i++ {
			ch <- x
			x, y = y, x+y
		}

		close(ch)
	}()
	// TODO: отправить последовательность Фибоначчи в канал
	return ch
}
