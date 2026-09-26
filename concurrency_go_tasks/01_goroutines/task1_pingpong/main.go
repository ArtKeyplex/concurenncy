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
	wg.Add(2)

	pingCh := make(chan string)
	pongCh := make(chan string)

	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			<-pingCh
			fmt.Fprintln(w, "ping")
			pongCh <- ""
		}

	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			<-pongCh
			fmt.Fprintln(w, "pong")

			if i < 4 {
				pingCh <- ""
			}
		}

	}()

	pingCh <- ""

	wg.Wait()

}

func main() {
	PingPong(os.Stdout)

}
