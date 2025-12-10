package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

var wg = &sync.WaitGroup{}

// PingPong должен запускать две горутины "ping" и "pong",
// которые поочередно выводят строки пять раз каждая.
// Реализуйте синхронизацию через каналы и ожидание завершения.
func PingPong(w io.Writer) {
	// TODO: реализовать обмен сообщениями между горутинами
	ping := make(chan struct{})
	pong := make(chan struct{})
	done := make(chan struct{})
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			<-ping
			fmt.Fprintf(w, "ping\n")
			pong <- struct{}{}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			<-pong

			fmt.Fprintf(w, "pong\n")
			switch {
			case i < 4:
				ping <- struct{}{}
			default:
				done <- struct{}{}
			}
		}
	}()

	ping <- struct{}{}
	<-done
}

func main() {
	PingPong(os.Stdout)
	wg.Wait()
}
