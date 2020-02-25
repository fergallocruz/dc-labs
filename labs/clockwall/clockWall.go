package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

//clockWall NewYork=localhost:8010 Tokyo=localhost:8020 London=localhost:8030
var timers []string

func main() {
	arguments := os.Args[1:]
	if len(arguments) < 1 {
		fmt.Println("Please provide a port number!")
		return
	}
	for _, s := range arguments {
		res1 := strings.Split(s, "=")
		res2 := strings.Split(res1[1], ":")
		PORT := ":" + res2[1]
		conn, err := net.Dial("tcp", PORT)
		println()
		if err != nil {
			if _, t := err.(*net.OpError); t {
				fmt.Println("Some problem connecting.")
			} else {
				fmt.Println("Unknown error: " + err.Error())
			}
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		scanner := bufio.NewScanner(conn)
		for {
			ok := scanner.Scan()
			text := scanner.Text()
			print(text)
			if !ok {
				fmt.Println("Reached EOF on server connection.")
				break
			}
		}
	}
}
