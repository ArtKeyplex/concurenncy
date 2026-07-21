package pipeline

// Run строит конвейер из трёх стадий: квадрат, умножение на 2 и суммирование.
func Run(nums []int) int {
	multiply := make(chan int)
	summ := make(chan int)
	ans := 0
	go func() {
		for _, v := range nums {
			multiply <- v * v
		}
		defer close(multiply)
	}()

	go func() {
		for rec := range multiply {
			summ <- rec * 2
		}
		defer close(summ)
	}()

	for rec := range summ {
		ans += rec
	}

	return ans
}
