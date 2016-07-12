package main

import (
	"encoding/json"
	"log"

	"github.com/robertojrojas/microservices-go/pets/pets-service/executor"
	"github.com/robertojrojas/microservices-go/pets/pets-service/services"
)

func main() {

	catService := &services.CatService{
		URL: "http://localhost:8091/api/cats",
	}
	birdsService := &services.BirdService{
		ServiceAddress: "localhost:8092",
	}
	dogsService := &services.DogsService{
		ServiceAddress: "amqp://guest:guest@localhost:5672/",
		RPCQueue:       "dog_service_rpc_queue",
	}
	rpcExecutor := executor.NewRPCExecutor(catService, birdsService, dogsService)
	results, err := rpcExecutor.GetAllPets()
	if err != nil {
		log.Fatal(err)
	}

	output, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("Unable to marshal output %s\n", err)
	}
	log.Printf("%s\n", output)

}
