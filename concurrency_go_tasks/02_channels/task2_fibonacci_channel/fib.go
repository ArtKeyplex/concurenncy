package fibonacci

// Fib возвращает канал, из которого можно читать первые n чисел Фибоначчи.
func Fib(n int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		if n <= 0 {
			return
		}

		x, y := 0, 1

		for i := 0; i < n; i++ {
			ch <- x
			x, y = y, y+x
		}
	}()

	return ch
}
