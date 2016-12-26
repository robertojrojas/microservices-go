package main

import (
	"flag"
	"log"

	"fmt"
	"github.com/robertojrojas/microservices-go/pets/dogs-service/messaging"
	"github.com/robertojrojas/microservices-go/pets/dogs-service/models"
)

var (
	uri               = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	queue             = flag.String("queue", "dog_service_rpc_queue", "Queue to get RPC messages from")
	mongoDBURI        = flag.String("mongoDBURI", "mongodb://localhost:27017", "MongoDB URI")
	mongoDBName       = flag.String("mongoDBName", "dogsDB", "MongoDB Database Name")
	mongoDBCollection = flag.String("mongoDBCollection", "dogs", "MongoDB Collection Name")
)

const moduleName = "dog-service"

var Version = "0.0.1"

var BuildTime string

func getVersion() string {
	return fmt.Sprintf("%s version: %s build time: %s", moduleName, Version, BuildTime)
}

func main() {

	flag.Parse()

	log.Printf("%s\n", getVersion())

	dogMongoStore := models.NewDogMongoStore(*mongoDBURI, *mongoDBName, *mongoDBCollection)
	defer dogMongoStore.Disconnect()

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
