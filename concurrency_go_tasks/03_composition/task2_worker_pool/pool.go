package pool

import (
	"sync"
)

// RunPool обрабатывает задачи параллельно в заданном количестве воркеров
// и возвращает сумму результатов.

// TODO: реализовать пул воркеров и сбор результатов
func RunPool(jobs []int, workers int) int {

	if workers <= 0 { // если воркеров нет или их количество отрицательно
		workers = 1 // должен работать 1 воркер
	}

	chjobs := make(chan int)           // канал для задач
	resch := make(chan int, len(jobs)) // буферизирваный канал для результатов

	var wg sync.WaitGroup // группа ожидания

	for i := 0; i < workers; i++ { // запускаем горутины-воркеры
		wg.Add(1)   // увеличиваем счётчик
		go func() { // запускаем воркера
			for num := range chjobs { // читаем задачи из канала chjobs
				result := num   //
				resch <- result // отправляем результаты в канал
			}
			wg.Done() // ждём пока все воркеры закончат работу
		}()

	}
	for _, job := range jobs {
		chjobs <- job // отправляем каждую задачу в канал
	}
	close(chjobs) // закрыли канал

	wg.Wait()    // ждём пока все закончат
	close(resch) // закрываем 2й канал

	total := 0 // создаём переменную для результатов

	for nums := range resch { // читаем результаты из resch
		total += nums // и складываем их
	}

	return total // возвращаем результат

}
