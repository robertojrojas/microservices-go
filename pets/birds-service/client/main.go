package main

import (
	"fmt"
	"log"
	"os"

	pb "github.com/robertojrojas/microservices-go/pets/birds-service/models"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8092"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBirdRepositoryClient(conn)

	if len(os.Args) > 1 {
		bird := &pb.Bird{
			Name: "chicken little",
			Age:  5,
			Type: pb.Bird_BLACKBIRDSCHICKADEE,
		}
		fmt.Println("Creating bird")
		bird, err = c.CreateBird(context.Background(), bird)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Contact the server
	// and print out its response.
	r, err := c.AllBirds(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for _, bird := range r.Birds {
		log.Printf("Results: %#v\n", bird)
	}

}
