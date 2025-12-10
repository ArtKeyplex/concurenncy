package pipeline

// Run строит конвейер из трёх стадий: квадрат, умножение на 2 и суммирование.
func Run(nums []int) int {
	// TODO: реализовать конвейер обработки чисел
	sqrt := Square(nums)
	mult := Multiply(sqrt)
	return Sum(mult)
}

func Square(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, v := range nums {
			out <- v * v
		}
	}()
	return out
}
func Multiply(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for v := range in {
			out <- v * 2
		}
	}()
	return out
}
func Sum(in <-chan int) int {
	total := 0
	for v := range in {
		total += v
	}
	return total
}
