package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

func PingPong(w io.Writer) {
	ping := make(chan struct{})
	pong := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			<-ping
			fmt.Fprintln(w, "ping")
			pong <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			<-pong
			fmt.Fprintln(w, "pong")
			if i < 4 { 
				ping <- struct{}{}
			}
		}
	}()

	ping <- struct{}{}
	wg.Wait()
}

func main() {
	PingPong(os.Stdout)
}