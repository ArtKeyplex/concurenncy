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
	pingCh := make(chan struct{})
	pongCh := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(pingCh)

		for i := 0; i < 5; i++ {
			<-pongCh
			fmt.Fprintln(w, "ping")
			pingCh <- struct{}{}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(pongCh)

		for i := 0; i < 5; i++ {
			pongCh <- struct{}{}
			<-pingCh
			fmt.Fprintln(w, "pong")
		}
	}()

	wg.Wait()
}

func main() {
	PingPong(os.Stdout)
}
