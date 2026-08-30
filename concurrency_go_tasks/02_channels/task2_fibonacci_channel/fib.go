package fibonacci

// Fib возвращает канал, из которого можно читать первые n чисел Фибоначчи.
func Fib(n int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)

		a := 0
		b := 1

		for range n {
			ch <- a

			next := a + b
			a = b
			b = next
		}
	}()
	return ch
}
