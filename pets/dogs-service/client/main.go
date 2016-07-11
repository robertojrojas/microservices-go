package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/robertojrojas/microservices-go/pets/dogs-service/messaging"
	"github.com/robertojrojas/microservices-go/pets/dogs-service/models"
	"github.com/streadway/amqp"
)

var (
	requestType = flag.Int("requestType", 0, "RPC Request Type")
	dogID       = flag.String("dogID", "", "Dog ID to query for")
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func readAllRPC() (err error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	rpcRequest := createRPCRequest()

	rpcData, err := json.Marshal(rpcRequest)
	failOnError(err, "Failed to Marshal rpcRequest")

	err = ch.Publish(
		"", // exchange
		"dog_service_rpc_queue", // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          rpcData,
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		rpcReply := &messaging.RPCMessage{}
		err := json.Unmarshal(d.Body, rpcReply)
		if err != nil {
			failOnError(err, "Failed to convert body to RPCMessage")
		}
		fmt.Printf("%s\n", string(rpcReply.Data))
		break
	}

	return
}

func createRPCRequest() (rpcRequest *messaging.RPCMessage) {

	switch messaging.RPCMessageType(*requestType) {
	case messaging.ReadAllMessage:
		rpcRequest = readAllDogsRequest()
	case messaging.CreateMessage:
		rpcRequest = createDogRequest()
	case messaging.ReadMessage:
		if *dogID == "" {
			failOnError(errors.New("-dogID flag is missing"), "")
		}
		rpcRequest = readDogRequest(*dogID)
	}

	return
}

func createDogRequest() *messaging.RPCMessage {
	dog := &models.Dog{
		Name: "MyDog",
		Age:  5,
		Type: "Super",
	}

	dogData, err := json.Marshal(dog)
	failOnError(err, "Failed to Marshal Dog object")

	rpcRequest := &messaging.RPCMessage{
		Type: messaging.CreateMessage,
		Data: dogData,
	}

	return rpcRequest
}

func readAllDogsRequest() *messaging.RPCMessage {
	rpcRequest := &messaging.RPCMessage{
		Type: messaging.ReadAllMessage,
		Data: []byte(""),
	}

	return rpcRequest
}

func readDogRequest(id string) *messaging.RPCMessage {
	rpcRequest := &messaging.RPCMessage{
		Type: messaging.ReadMessage,
		Data: []byte(id),
	}

	return rpcRequest
}

func main() {
	flag.Parse()

	rand.Seed(time.Now().UTC().UnixNano())

	log.Println(" [x] Requesting readAllRPC()")
	err := readAllRPC()
	failOnError(err, "Failed to handle RPC request")

}
