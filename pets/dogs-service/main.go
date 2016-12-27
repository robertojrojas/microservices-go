package main

import (
	"fmt"
	"github.com/robertojrojas/microservices-go/pets/dogs-service/messaging"
	"github.com/robertojrojas/microservices-go/pets/dogs-service/models"
	"log"
	"os"
)

const moduleName = "dog-service"

var Version = "0.0.1"

var BuildTime string

func getVersion() string {
	return fmt.Sprintf("%s version: %s build time: %s", moduleName, Version, BuildTime)
}

type config struct {
	rabbitMQURI       string
	rabbitMQQueue     string
	mongoDBURI        string
	mongoDBName       string
	mongoDBCollection string
}

func main() {

	log.Printf("%s\n", getVersion())

	envConfig := getConfig()
	log.Printf("%#v\n", envConfig)

	dogMongoStore, err := models.NewDogMongoStore(envConfig.mongoDBURI, envConfig.mongoDBName, envConfig.mongoDBCollection)
	checkErr(err, "connecting to MongoDB")
	defer dogMongoStore.Disconnect()

	messageHanlder := messaging.NewMessageHandler(dogMongoStore)
	amqpManager := messaging.NewAMQManager(envConfig.rabbitMQURI, envConfig.rabbitMQQueue, messageHanlder)
	err = amqpManager.Connect()
	checkErr(err, "connecting to RabbitMQ")

	defer amqpManager.Disconnect()

	log.Println("Ready to receive messages...")
	err = amqpManager.WaitForMessages()
	checkErr(err, "receive messages from RabbitMQ")

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

	envConfig.mongoDBURI = os.Getenv("MONGODB_URI")
	if envConfig.mongoDBURI == "" {
		envConfig.mongoDBURI = "mongodb://localhost:27017"
	}

	envConfig.mongoDBName = os.Getenv("MONGODB_DBNAME")
	if envConfig.mongoDBName == "" {
		envConfig.mongoDBName = "dogsDB"
	}

	envConfig.mongoDBCollection = os.Getenv("MONGODB_COLLECTION")
	if envConfig.mongoDBCollection == "" {
		envConfig.mongoDBCollection = "dogs"
	}

	return &envConfig
}

func checkErr(err error, message string) {
	if err != nil {
		log.Fatalf("error: %s - %s", message, err)
	}

}
