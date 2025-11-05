package main

import (
	"io"
	"os"
	"sync"
)

// PingPong должен запускать две горутины "ping" и "pong",
// которые поочередно выводят строки пять раз каждая.
// Реализуйте синхронизацию через каналы и ожидание завершения.
func PingPong(w io.Writer) {
	// TODO: реализовать обмен сообщениями между горутинами
	ch := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			w.Write([]byte("ping\n"))
			ch <- 1
			<-ch
		}
		close(ch)

	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			<-ch
			w.Write([]byte("pong\n"))
			ch <- 1
		}
	}()
	wg.Wait()
}

func main() {
	PingPong(os.Stdout)
}
