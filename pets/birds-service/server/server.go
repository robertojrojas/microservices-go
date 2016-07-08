package server

import (
	"flag"
	"fmt"
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

func (s *server) AllBirds(ctx context.Context, in *pb.Empty) (birdCatalog *pb.BirdCatalog, err error) {

	birdRecords, err := s.birdsDBStore.ReadAllBirds()
	if err != nil {
		return
	}

	birdCatalog = &pb.BirdCatalog{}
	for _, birdRecord := range birdRecords {
		bird := &pb.Bird{
			Name: birdRecord.Name,
			Id:   birdRecord.ID,
			Age:  birdRecord.Age,
			Type: pb.BirdTypeFromString(birdRecord.Type),
		}
		birdCatalog.Birds = append(birdCatalog.Birds, bird)
	}

	return

}

func (s *server) CreateBird(ctx context.Context, in *pb.Bird) (out *pb.Bird, err error) {

	birdRecord := &pb.BirdRecord{
		Name: in.Name,
		Age:  in.Age,
		Type: pb.GetBirdTypeStringFromType(in.Type),
	}

	err = s.birdsDBStore.CreateBird(birdRecord)
	if err != nil {
		fmt.Println("Err", err)
		return
	}

	out = &pb.Bird{
		Id:   birdRecord.ID,
		Name: birdRecord.Name,
		Age:  birdRecord.Age,
		Type: pb.BirdTypeFromString(birdRecord.Type),
	}

	return

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
