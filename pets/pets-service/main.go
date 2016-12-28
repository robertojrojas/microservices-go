package main

import (
	"encoding/json"
	"log"

	"github.com/robertojrojas/microservices-go/pets/pets-service/executor"
	"github.com/robertojrojas/microservices-go/pets/pets-service/services"
	"os"
)

type config struct {
	catsServiceURI string
	birdsServiceURI string
	rabbitMQURI   string
	rabbitMQQueue string
}

func main() {

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
