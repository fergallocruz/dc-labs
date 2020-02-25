// Clock2 is a concurrent TCP server that periodically writes the time.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var port int
var tz string

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		t, err := TimeIn(time.Now(), tz)
		if err == nil {
			loc := t.Location()
			fmt.Println(loc, t.Format("\t:15:04:05"))
		} else {
			fmt.Println(tz, "")
		}
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}
func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

func main() {
	tz = os.Getenv("TZ")
	port := flag.Int("port", port, "test port")
	// Parse the flags.
	flag.Parse()
	fmt.Println("  Port:", *port)
	fmt.Println("  TZ:", tz)
	listener, err := net.Listen("tcp", "localhost:"+(strconv.Itoa(*port)))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle connections concurrently
	}
}

//TZ=Europe/London ./clock2 -port 8030 &
