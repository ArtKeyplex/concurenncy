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
	pingChan := make(chan struct{})
	pongChan := make(chan struct{})
	wg := &sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			w.Write([]byte("ping\n"))
			pongChan <- struct{}{}
			<-pingChan
		}
	}()

	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			<-pongChan
			w.Write([]byte("pong\n"))
			pingChan <- struct{}{}
		}
	}()

	wg.Wait()
}

func main() {
	PingPong(os.Stdout)
}
