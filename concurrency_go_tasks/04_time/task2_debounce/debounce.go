package debounce

import "time"

// Debounce принимает значения и отдаёт только последнее после паузы d.
// TODO: реализовать дебаунс значений из канала
func Debounce(d time.Duration, in <-chan int) <-chan int {
	ch := make(chan int)

	go func() {
		for {

			num, ok := <-in // Ждём первое число новой группы.

			if !ok { // Если вход закрыт — закрываем выход и заканчиваем.
				close(ch)
				return
			}

			last := num // Запоминаем последнее число.

			timer := time.NewTimer(d) // Создаём таймер на паузу d.

		wait:
			for {
				select { // Ждём либо новое число, либо конец паузы.
				case num, ok := <-in:

					if !ok {
						<-timer.C
						ch <- last
						close(ch)
						return
					}

					last = num     // Пришло новое число: оно становится последним,
					timer.Reset(d) // а таймер перезапускается.

				case <-timer.C:
					ch <- last // Пауза прошла: отправляем последнее число.
					break wait
				}
			}
		}
	}()

	return ch
}
