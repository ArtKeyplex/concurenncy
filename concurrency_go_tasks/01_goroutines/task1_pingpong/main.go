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

	ch := make(chan int)
	ping := func() {
		fmt.Fprintln(w, "ping")
		ch <- 1
		wg.Done()
	}
	pong := func() {
		<-ch
		fmt.Fprintln(w, "pong")
		wg.Done()
	}

	for i := 0; i < 5; i++ {
		wg.Add(2)
		go ping()
		go pong()
		wg.Wait()
	}
}

func main() {
	PingPong(os.Stdout)
}
