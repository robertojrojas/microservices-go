package server

import (
	"encoding/json"
	"log"

	"fmt"
	"github.com/gorilla/mux"
	"github.com/robertojrojas/microservices-go/pets/pets-service/executor"
	"github.com/robertojrojas/microservices-go/pets/pets-service/services"
	"net/http"
	"os"
)

type PetsServer struct {
	rpcExecutor *executor.RPCExecutor
}

type config struct {
	catsServiceURI  string
	birdsServiceURI string
	rabbitMQURI     string
	rabbitMQQueue   string
}

// StartServer configures and starts API Server
func StartServer(serverHostPort string) error {

	appConfig := getConfig()

	catService := &services.CatService{
		URL: appConfig.catsServiceURI,
	}
	birdsService := &services.BirdService{
		ServiceAddress: appConfig.birdsServiceURI,
	}

	dogsService := &services.DogsService{
		ServiceAddress: appConfig.rabbitMQURI,
		RPCQueue:       appConfig.rabbitMQQueue,
	}
	rpcExecutor := executor.NewRPCExecutor(catService, birdsService, dogsService)

	petsServer := &PetsServer{
		rpcExecutor: rpcExecutor,
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/pets", petsServer.petsHandler).Methods("GET")

	http.Handle("/", router)

	log.Printf("Listening on [%s]....\n", serverHostPort)
	return http.ListenAndServe(serverHostPort, nil)
}

func (s *PetsServer) petsHandler(w http.ResponseWriter, r *http.Request) {
	results, err := s.rpcExecutor.GetAllPets()
	if err != nil {
		log.Fatal(err)
	}

	output, err := json.Marshal(results)
	if err != nil {
		fmt.Fprintf(w, "Unable to marshal output %s\n", err)
		return
	}

	fmt.Fprintf(w, "%s\n", output)
}

func getConfig() *config {

	envConfig := config{}

	envConfig.rabbitMQURI = os.Getenv("RABBITMQ_URI")
	if envConfig.rabbitMQURI == "" {
		envConfig.rabbitMQURI = "amqp://guest:guest@localhost:5672/"
	}

	envConfig.rabbitMQQueue = os.Getenv("RABBITMQ_QUEUE")
	if envConfig.rabbitMQQueue == "" {
		envConfig.rabbitMQQueue = "dog_service_rpc_queue"
	}

	envConfig.catsServiceURI = os.Getenv("CATS_SERVICE_URI")
	if envConfig.catsServiceURI == "" {
		envConfig.catsServiceURI = "http://cats-service-svc:8091/api/cats"
	}

	envConfig.birdsServiceURI = os.Getenv("BIRDS_SERVICE_URI")
	if envConfig.birdsServiceURI == "" {
		envConfig.birdsServiceURI = "birds-service-svc:8092"
	}

	return &envConfig
}
