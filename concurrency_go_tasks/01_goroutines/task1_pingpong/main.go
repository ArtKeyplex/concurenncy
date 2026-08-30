package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

func PingPong(w io.Writer) {
	wg := sync.WaitGroup{}
	ch := make(chan string)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			ch <- "ping"
			fmt.Fprintln(w, "ping", i+1)
			<-ch

		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			<-ch
			fmt.Fprintln(w, "pong", i+1)
			ch <- "pong"
		}
	}()

	wg.Wait()
}

func main() {
	PingPong(os.Stdout)
}
