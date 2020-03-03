package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	_, err := os.Create("file.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	t := make(chan time.Time)
	end := t

	for counter := 1; ; counter++ {
		go read(end, counter)
		t <- time.Now()

		newStart := end
		end = make(chan time.Time)
		go connect(newStart, end)
	}
}

func read(end chan time.Time, counter int) {
	startTime := <-end
	endTime := time.Now()

	d := fmt.Sprintf("goroutine: %d time: %v", counter, endTime.Sub(startTime))
	file, err := os.OpenFile("file.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(d)
}

func connect(src, dst chan time.Time) {
	for t := range src {
		dst <- t
	}
}
