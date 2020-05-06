package controller

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.nanomsg.org/mangos"
	"go.nanomsg.org/mangos/protocol/pub"
	"google.golang.org/grpc/metadata"

	// register transports
	_ "go.nanomsg.org/mangos/transport/all"
)

var controllerAddress = "tcp://localhost:40899"
var sock mangos.Socket

type MD map[string][]string

var Nodes = make(map[string]metadata.MD)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}

func Start() {
	var err error

	if sock, err = pub.NewSocket(); err != nil {
		die("can't get new pub socket: %s", err)
	}
	if err = sock.Listen(controllerAddress); err != nil {
		die("can't listen on pub socket: %s", err.Error())
	}
	defer sock.Close()
	for {

		// Could also use sock.RecvMsg to get header
		fmt.Printf(sock.Info().PeerName)
		d := date()
		log.Printf("Controller: Publishing Date %s\n", d)
		if err = sock.Send([]byte(d)); err != nil {
			die("Failed publishing: %s", err.Error())
		}
		time.Sleep(time.Second * 3)
	}

}

/*func Register(workerName string, tags string, status string, usage int) {
	Nodes[workerName] = Node{
		worker: workerName,
		tags:   tags,
		status: "im fine",
		usage:  70,
	}
}*/
