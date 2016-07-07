package main

import (
	"log"

	pb "github.com/robertojrojas/microservices-go/pets/birds-service/models"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBirdRepositoryClient(conn)

	// Contact the server and print out its response.

	r, err := c.AllBirds(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Results: %#v\n", r.Message)
}
