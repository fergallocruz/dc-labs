package main

import (
        "fmt"
        "os"
        "os/signal"
        "runtime/trace"
)

func routine(id int) {
        fmt.Printf("Starting #%v goroutine\n", id)
        //for {
        //}
}

func main() {

        // Enable tracing
        f, err := os.Create("trace.out")
        if err != nil {
                panic(err)
        }
        defer f.Close()

        err = trace.Start(f)
        if err != nil {
                panic(err)
        }
        defer trace.Stop()

        // Create Signal to catch the Ctrl + C
        c := make(chan os.Signal, 1)
        signal.Notify(c)

        fmt.Println("Start Main Goroutine")

        fmt.Println("Starting new Goroutines")
        for i := 0; i < 100; i++ {
                go routine(i)
        }
        s := <-c
        fmt.Println("Got signal to stop: ", s)
}
