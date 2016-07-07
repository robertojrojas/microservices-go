package server

import (
	"flag"
	"log"
	"net"

	pb "github.com/robertojrojas/microservices-go/pets/birds-service/models"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// server is used to implement birds-service Server.
type server struct {
	birdsDBStore pb.BirdsDataStore
}

var serverHostPort string
var dbURL *string

func init() {

	dbURL = flag.String("dbURL", "user=birdman password=mycape dbname=birddb sslmode=disable", "DB URL for PostgreSQL")
	flag.StringVar(&serverHostPort, "serverHostPort", ":8092", "Host and port server listens on")

}

// NewBirdsServer return a new server
func NewBirdsServer(dataStore pb.BirdsDataStore) *server {
	return &server{
		birdsDBStore: dataStore,
	}
}

func (s *server) AllBirds(ctx context.Context, in *pb.Empty) (*pb.BirdCatalog, error) {

	birdCatalog := &pb.BirdCatalog{}

	return birdCatalog, nil

}

func (s *server) CreateBird(ctx context.Context, in *pb.Bird) (*pb.Bird, error) {

	bird := &pb.Bird{}

	return bird, nil

}

func (s *server) ReadBird(ctx context.Context, in *pb.BirdId) (*pb.Bird, error) {

	bird := &pb.Bird{}

	return bird, nil

}

func (s *server) UpdateBird(ctx context.Context, in *pb.Bird) (*pb.Bird, error) {

	bird := &pb.Bird{}

	return bird, nil

}

func (s *server) DeleteBird(ctx context.Context, in *pb.BirdId) (*pb.Empty, error) {

	return &pb.Empty{}, nil

}

// StartServer configures and starts API Server
func StartServer() error {

	birdsDB := pb.NewBirdsDB(*dbURL)
	birdServer := NewBirdsServer(birdsDB)

	lis, err := net.Listen("tcp", serverHostPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBirdRepositoryServer(s, birdServer)
	return s.Serve(lis)

}
