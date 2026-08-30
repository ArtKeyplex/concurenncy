package producerconsumer

import (
	"io"
	"strconv"
	"sync"
)

// Run запускает продюсера, который отправляет числа от 1 до 10, и консюмера,
// который выводит их в writer. Используйте небуферизованный канал и ожидание
// завершения горутин.
func Run(w io.Writer) {
	intCh := make(chan int)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer close(intCh)

		for i := 1; i <= 10; i++ {
			intCh <- i
		}
	}()

	go func() {
		defer wg.Done()
		for i := range intCh {
			w.Write([]byte(strconv.Itoa(i) + "\n"))
		}
	}()

	wg.Wait()
}
