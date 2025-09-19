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

	pingCh := make(chan struct{}, 1)
	pongCh := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			<-pingCh
			fmt.Fprintln(w, "ping")
			pongCh <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			<-pongCh
			fmt.Fprintln(w, "pong")
			pingCh <- struct{}{}
		}
	}()

	pingCh <- struct{}{}

	wg.Wait()
}

func main() {
	PingPong(os.Stdout)
}
