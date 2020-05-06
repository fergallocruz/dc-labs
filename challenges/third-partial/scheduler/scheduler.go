package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/CodersSquad/dc-labs/challenges/third-partial/controller"
	pb "github.com/CodersSquad/dc-labs/challenges/third-partial/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//const (
//	address     = "localhost:50051"
//	defaultName = "world"
//)

type Job struct {
	Address string
	RPCName string
}

func schedule(job Job) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(job.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	controller.Nodes["Popeye"] = metadata.New(map[string]string{"controller": "", "host": "", "node": "", "tags": "", "status": ""})
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: job.RPCName})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Scheduler: RPC respose from %s : %s", job.Address, r.GetMessage())

}

func Start(jobs chan Job) error {
	for {
		job := <-jobs
		schedule(job)
	}
	return nil
}
