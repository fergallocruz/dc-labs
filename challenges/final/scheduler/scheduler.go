package scheduler

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/CodersSquad/dc-labs/challenges/final/controller"

	pb "github.com/CodersSquad/dc-labs/challenges/final/proto"
	"google.golang.org/grpc"
)

//const (
//	address     = "localhost:50051"
//	defaultName = "world"
//)

type Job struct {
	Address string
	RPCName string
}

var counter int

func schedule(job Job, name string) {

	// Set up a connection to the server.
	conn, err := grpc.Dial(job.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	controller.ChangeStatus(name)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: job.RPCName})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Scheduler: RPC respose from %s : %s", job.Address, r.GetMessage())
	controller.ChangeStatus(name)
	counter++
}

func Start(jobs chan Job) error {
	counter = 0
	for {
		job := <-jobs
		time.Sleep(time.Second * 5)
		lowestUsage := 99999
		lowestPort := 0
		worker := controller.Worker{}
		for _, data := range controller.Nodes {
			if data.Usage < lowestUsage {
				lowestPort = data.Port
				lowestUsage = data.Usage
				worker = data
			}
		}
		controller.IncreaseUse(worker.Name)
		controller.Register(worker.Name, counter)
		if lowestPort == 0 {
			return nil
		}

		job.Address = "localhost:" + strconv.Itoa(lowestPort)
		schedule(job, worker.Name)
	}
	return nil
}