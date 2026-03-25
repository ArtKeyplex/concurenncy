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
	pingCh := make(chan struct{})
	pongCh := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for i:=0; i<5; i++ {
			<-pingCh
			w.Write([]byte("ping\n"))
			pongCh <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()

		for i:=0; i<5; i++ {
			<-pongCh
			w.Write([]byte("pong\n"))
			if i<4 {
				pingCh<-struct{}{}
			}
		}
	}()

	pingCh <- struct{}{}
	wg.Wait()
}

func main() {
	PingPong(os.Stdout)
}