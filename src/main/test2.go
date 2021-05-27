package main

import (
	"context"
	"log"
	"time"
	"google.golang.org/grpc"
	"github.com/JijiaZan/godml/pyserver"
	"gonum.org/v1/gonum/mat"
)

const (
	address     = "localhost:1888"

)

func main() {
	data := []float64{ 3, 4, 5, 5, 2, 4, 7, 8, 11, 8, 12,
		11, 13, 13, 16, 17, 18, 17, 19, 21}
	y2 := mat.NewDense(4, 5, data)


	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	log.Printf("ok")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pyserver.NewWorkerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Push(ctx, &pyserver.PushRequest{
		T: 2,
		Weight: []float32{1,2,3},
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetSuccess())
}