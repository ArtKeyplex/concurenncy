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
	var wg = &sync.WaitGroup{}
	// TODO: реализовать продюсер и консюмер
	ch := make(chan int)
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		close(ch)

	}()
	go func() {
		defer wg.Done()

		for val := range ch {
			fmt.Fprintf(w, "%d\n", val)
		}

	}()

	wg.Wait()
}
