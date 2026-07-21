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
	ch := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i < 11; i++ {
			ch <- i
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 10 {
			fmt.Fprintln(w, <-ch)
		}
	}()

	wg.Wait()

}
