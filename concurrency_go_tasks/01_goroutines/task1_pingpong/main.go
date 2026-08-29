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
	wg.Add(2)

	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			<-pingCh

			if _, err := w.Write([]byte("ping\n")); err != nil {
				fmt.Println(err)
				return
			}

			pongCh <- struct{}{}
		}
	}()
	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			<-pongCh

			if _, err := w.Write([]byte("pong\n")); err != nil {
				fmt.Println(err)
				return
			}

			if i < 4 {
				pingCh <- struct{}{}
			}
		}
	}()

	pingCh <- struct{}{}

	wg.Wait()
}

func main() {
	PingPong(os.Stdout)
}
