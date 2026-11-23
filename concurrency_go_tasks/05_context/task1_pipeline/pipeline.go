package pipelinectx

import (
	"context"
	"sync"
)

// Run строит конвейер из двух стадий: удвоение и суммирование.
// Конвейер должен останавливаться, если ctx отменён.
// Возвращает итоговую сумму и ошибку контекста при отмене.

// TODO: реализовать конвейер с остановкой по ctx
func Run(ctx context.Context, nums []int) (int, error) {
	if err := ctx.Err(); err != nil { // если контекст уже отменён — сразу возвращаем ошибку
		return 0, err
	}
	var wg sync.WaitGroup     // нужен, чтобы дождаться завершения обеих горутин
	result := make(chan int)  // канал для суммы
	doubled := make(chan int) // канал для удвоеных чисел

	wg.Add(2) // сообщаем что будет 2 горутины

	go func() { // горутина удвоения
		defer wg.Done()            // wg.Done() вызовется один раз
		for _, num := range nums { // проходим по всем числам
			select {
			case <-ctx.Done(): // контекст отменён — останавливаемся
				return
			case doubled <- num * 2: // отправили удвоенное число
			}
		}
		close(doubled) // закрыли канал
	}()

	go func() { // горутина суммы
		defer wg.Done()
		defer close(result) // когда горутина закончится - закрытваем канал

		total := 0 // тут будет итоговая сумма

		for {
			select {
			case <-ctx.Done(): // контекст отменён - стоп
				return

			case num, ok := <-doubled: // если канал doubled закрыт — все числа обработаны
				if !ok {
					result <- total // отправляем итоговую сумму
					return
				}

				total += num // складываем удвоенное число
			}
		}
	}()

	select { // ждём либо результат, либо отмену контекста
	case total := <-result:
		wg.Wait()         // дожидаемся завершения горутин
		return total, nil // возвращаем сумму

	case <-ctx.Done():
		wg.Wait()           // дожидаемся завершения горутин
		return 0, ctx.Err() // возвращаем ошибку контекста
	}
}
