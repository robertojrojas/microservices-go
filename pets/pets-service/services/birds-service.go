package services

import (
	"log"

	"google.golang.org/grpc"

	pb "github.com/robertojrojas/microservices-go/pets/birds-service/models"
)

const (
	BirdServiceKey = "BirdService"
)

type BirdService struct {
	ServiceAddress string
}

func (service *BirdService) RPC(rpcRequest *RPCRequest) (rpcResponse *RPCResponse, err error) {

	// Set up a connection to the server.
	conn, err := grpc.Dial(service.ServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBirdRepositoryClient(conn)

	r, err := c.AllBirds(rpcRequest.Ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	rpcResponse = &RPCResponse{
		Key: BirdServiceKey,
	}
	rpcResponse.Data = r.Birds

	return

}
