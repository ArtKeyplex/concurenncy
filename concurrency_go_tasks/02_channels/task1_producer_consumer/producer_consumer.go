package producerconsumer

import (
	"fmt"
	"io"
	"sync"
)

// Run запускает продюсера, который отправляет числа от 1 до 10, и консюмера,
// который выводит их в writer. Используйте небуферизованный канал и ожидание
// завершения горутин.
func Run(w io.Writer) {
	// TODO: реализовать продюсер и консюмер
	ch := make(chan int)
	var wg sync.WaitGroup

	producer := func() {

		for i := 1; i <= 10; i++ {
			ch <- i
		}
		wg.Done()
	}

	consumer := func() {
		for i := 1; i <= 10; i++ {
			fmt.Fprintln(w, <-ch)
		}

		wg.Done()
	}

	wg.Add(2)

	go producer()
	go consumer()

	wg.Wait()
}
