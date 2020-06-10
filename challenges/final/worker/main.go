package main

import (
	"context"
	"dc-labs/mangos/protocol/respondent"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "github.com/CodersSquad/dc-labs/challenges/final/proto"
	"go.nanomsg.org/mangos"
	"google.golang.org/grpc"

	// register transports
	_ "go.nanomsg.org/mangos/transport/all"
)

var (
	defaultRPCPort = 50051
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

var (
	controllerAddress = ""
	workerName        = ""
	tags              = ""
	status            = ""
	workDone          = 0
	usage             = 0
	port              = 0
	jobsDone          = 0
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func init() {
	flag.StringVar(&controllerAddress, "controller", "tcp://localhost:40899", "Controller address")
	flag.StringVar(&workerName, "node-name", "hard-worker", "Worker Name")
	flag.StringVar(&tags, "tags", "gpu,superCPU,largeMemory", "Comma-separated worker tags")
}
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	switch in.Name {
	case "test":
		workDone++
		log.Printf("RPC [Worker] %+v: testing...", workerName)
		usage++
		status = "Running"
		usage--
		return &pb.HelloReply{Message: "Here " + workerName + " testing..."}, nil
	default:
		workDone++
		log.Printf("[Worker] %+v: calling", workerName)
		usage++
		status = "Running"
		return &pb.HelloReply{Message: "Hello " + workerName}, nil
	}
}

// joinCluster is meant to join the controller message-passing server
func joinCluster() {
	var sock mangos.Socket
	var err error
	var msg []byte
	if sock, err = respondent.NewSocket(); err != nil {
		die("can't get new respondent socket: %s", err.Error())
	}
	log.Printf("Connecting to controller on: %s", controllerAddress)
	if err = sock.Dial(controllerAddress); err != nil {
		die("can't dial on respondent socket: %s", err.Error())
	}
	for {
		if msg, err = sock.Recv(); err != nil {
			die("Cannot recv: %s", err.Error())
		}
		data := workerName + " " + status + " " + strconv.Itoa(usage) + " " + tags + " " + strconv.Itoa(defaultRPCPort) + " " + strconv.Itoa(jobsDone)
		if err = sock.Send([]byte(data)); err != nil {
			die("Cannot send: %s", err.Error())
		}

		log.Printf("Message-Passing: Worker(%s): Received %s\n", workerName, string(msg))
	}
}

func getAvailablePort() int {
	port := defaultRPCPort
	for {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
		if err != nil {
			port = port + 1
			continue
		}
		ln.Close()
		break
	}
	return port
}

func main() {
	flag.Parse()
	//jobsDone := 0

	// Subscribe to Controller
	go joinCluster()

	// Setup Worker RPC Server
	rpcPort := getAvailablePort()
	defaultRPCPort = rpcPort
	//data := workerName + "|" + status + "|" + strconv.Itoa(usage) + "|" + tags + "|" + strconv.Itoa(rpcPort) + "|" + strconv.Itoa(jobsDone)
	log.Printf("Starting RPC Service on localhost:%v", rpcPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", rpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}