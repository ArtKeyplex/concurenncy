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
	pingChan1 := make(chan struct{})
	pingChan2 := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 5 {
			fmt.Fprintln(w, "ping")
			pingChan2 <- struct{}{}
			<-pingChan1
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 5 {
			<-pingChan2
			fmt.Fprintln(w, "pong")
			pingChan1 <- struct{}{}
		}
	}()
	wg.Wait()
}
func main() {
	PingPong(os.Stdout)
}
