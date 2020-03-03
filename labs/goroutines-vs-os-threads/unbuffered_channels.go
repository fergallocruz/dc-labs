package main

import (
	"fmt"
	"time"
)

func main() {
	ping, pong := make(chan string), make(chan string)
	var counter float32

	t := time.NewTimer(1 * time.Minute)
	end := make(chan struct{})
	shutdown := make(chan struct{})

	// ping->pong
	go func() { //ping receives
	loop:
		for {
			select {
			case <-shutdown:
				break loop
			case v := <-ping:
				counter++
				v = "pong"
				fmt.Printf("\ro............")
				pong <- v //sends to pong
			}
		}
		end <- struct{}{}
	}()

	go func() { //pong receives
	loop:
		for {
			select {
			case <-shutdown:
				break loop
			case v := <-pong:
				v = "ping"
				fmt.Printf("\r............o")
				ping <- v
			}
		}
		end <- struct{}{}
	}()
	ping <- "ping"

	// 1 minute clock
	<-t.C
	close(shutdown)
	t.Stop()
	// gets last communication from the one on which it was stopped
	select {
	case <-ping:
	case <-pong:
	}
	// waits for ping and pong to end
	<-end
	<-end
	fmt.Printf("\nCOMUNICATIONS PER SECOND: %f ", counter/60)
}
