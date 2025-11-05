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
	// TODO: реализовать продюсер и консюмер
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i < 11; i++ {
			ch <- i
		}
		close(ch)
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case res, ok := <-ch:
				if ok {
					w.Write([]byte(strconv.Itoa(res) + "\n"))

				} else {
					return
				}
			}
		}
	}()
	wg.Wait()
}
