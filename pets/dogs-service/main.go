package main

import (
	"flag"
	"log"

	"github.com/robertojrojas/microservices-go/pets/dogs-service/messaging"
	"github.com/robertojrojas/microservices-go/pets/dogs-service/models"
)

var (
	uri   = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	queue = flag.String("queue", "dog_service_rpc_queue", "Queue to get RPC messages from")
)

func main() {
	flag.Parse()

	dogMongoStore := models.NewDogMongoStore()
	messageHanlder := messaging.NewMessageHandler(dogMongoStore)
	amqpManager := messaging.NewAMQManager(*uri, *queue, messageHanlder)
	err := amqpManager.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer amqpManager.Disconnect()

	err = amqpManager.WaitForMessages()
	if err != nil {
		log.Fatal(err)
	}
}
