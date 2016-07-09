package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/robertojrojas/microservices-go/pets/dogs-service/messaging"
	"github.com/streadway/amqp"
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

	rpcRequest := &messaging.RPCMessage{
		Type: messaging.ReadAllMessage,
		Data: []byte(""),
	}

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
		// if corrId == d.CorrelationId {
		// 	res, err = strconv.Atoi(string(d.Body))
		// 	failOnError(err, "Failed to convert body to integer")
		// 	break
		// }
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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	log.Println(" [x] Requesting readAllRPC()")
	err := readAllRPC()
	failOnError(err, "Failed to handle RPC request")

}
