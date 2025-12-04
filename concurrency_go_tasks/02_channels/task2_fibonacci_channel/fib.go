package fibonacci

// Fib возвращает канал, из которого можно читать первые n чисел Фибоначчи.
func Fib(n int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)

		a, b := 0, 1

		for i := 0; i < n; i++ {
			ch <- a

			a, b = b, a+b
		}
	}()

	return ch
}
