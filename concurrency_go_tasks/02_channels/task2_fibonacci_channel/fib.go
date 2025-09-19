package fibonacci

// Fib возвращает канал, из которого можно читать первые n чисел Фибоначчи.
func Fib(n int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		if n <= 0 {
			return
		}
		a, b := 0, 1
		if n >= 1 {
			ch <- a
		}
		for i := 1; i < n; i++ {
			ch <- b

			a, b = b, b+a
		}

	}()
	return ch
}
