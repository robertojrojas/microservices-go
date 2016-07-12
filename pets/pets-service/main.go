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
	rpcExecutor := executor.NewRPCExecutor(catService)
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
