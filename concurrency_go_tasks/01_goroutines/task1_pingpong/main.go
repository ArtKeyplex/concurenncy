package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// PingPong должен запускать две горутины "ping" и "pong",
// которые поочередно выводят строки пять раз каждая.
// Реализуйте синхронизацию через каналы и ожидание завершения.
func PingPong(w io.Writer) {
	// TODO: реализовать обмен сообщениями между горутинами
	var wg sync.WaitGroup

	ch := make(chan string)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			ch <- "ping"
			msg := <-ch
			fmt.Fprintln(w, msg)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			msg := <-ch
			fmt.Fprintln(w, msg)
			ch <- "pong"
		}
	}()

	wg.Wait()
}

func main() {
	PingPong(os.Stdout)
}
