package pipeline

// Run строит конвейер из трёх стадий: квадрат, умножение на 2 и суммирование.

// TODO: реализовать конвейер обработки чисел
func Run(nums []int) int {

	in := make(chan int) // создали канал in для исходных чисел
	go func() {          // запускаем горутину
		for _, num := range nums { // перебирам числа из nums
			in <- num // кладём текущее число в канал
		}
		close(in) // закрываем канал
	}()

	squared := make(chan int) // канал для квадратов

	go func() {
		for num := range in { // читаем числа из in
			squared <- num * num // кладём в канал квадрат чисел
		}
		close(squared) // закрываем канал
	}()

	doubled := make(chan int) // канал для квадратов умноженых на 2
	go func() {
		for num := range squared { // читаем числа из squared
			result := num * 2 // умножаем текущее число на 2
			doubled <- result // кладём результат в канал
		}
		close(doubled) // закрываем канал
	}()

	total := 0 // заводим переменную для суммы результатов

	for nums := range doubled { // читаем из doubled
		total += nums // прибавляем текущее число к общей сумме
	}
	return total // возвращаем сумму
}
