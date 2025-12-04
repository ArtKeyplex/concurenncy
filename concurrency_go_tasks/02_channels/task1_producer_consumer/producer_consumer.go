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
	var wg sync.WaitGroup

	numChan := make(chan int)

	wg.Add(2)

	go producer(1, 10, numChan, &wg)
	go consumer(w, numChan, &wg)

	wg.Wait()
}

func producer(left, right int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch)

	for i := left; i <= right; i++ {
		ch <- i
	}
}

func consumer(w io.Writer, producerChan chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := range producerChan {
		w.Write([]byte(fmt.Sprintf("%d\n", i)))
	}
}
