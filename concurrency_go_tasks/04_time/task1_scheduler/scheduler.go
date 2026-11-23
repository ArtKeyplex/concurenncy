package scheduler

import (
	"sync"
	"time"
)

// Every запускает f каждые d и возвращает функцию для остановки.
// TODO: периодический вызов функции с возможностью остановки

func Every(d time.Duration, f func()) (stop func()) {
	var once sync.Once // создали once, чтобы close(done) выполнился один раз

	ticker := time.NewTicker(d) // создали тикер с интарвалом d
	done := make(chan struct{}) // канал для остановки

	go func() { // запустили горутину
		for { // бесконечный цикл
			select { // ждём одно из 2х событий
			case <-ticker.C: // пришел тик, прошло d времени
				f() // вызываем функцию f()
			case <-done: // пришел сигнал остановки
				ticker.Stop() // останавливаем тикер
				return        // выходим из горутины
			}
		}
	}()

	stop = func() { // создали функцию stop
		once.Do(func() { // close(done) выполнится только один раз
			close(done)
		})
	}

	return stop // возвращаем stop
}
