package services

import (
	"encoding/json"
	"log"
	"math/rand"

	"github.com/robertojrojas/microservices-go/pets/dogs-service/messaging"
	"github.com/robertojrojas/microservices-go/pets/dogs-service/models"
	"github.com/streadway/amqp"
)

const (
	DogsServiceKey = "DogsService"
)

type DogsService struct {
	ServiceAddress string
	RPCQueue       string
}

func (service *DogsService) RPC(rpcRequest *RPCRequest) (rpcResponse *RPCResponse, err error) {

	log.Printf("[%T] Using for ServiceAddress %s RPCQueue %s\n", service, service.ServiceAddress, service.RPCQueue)
	conn, err := amqp.Dial(service.ServiceAddress)
	if err != nil {
		return
	}
	defer conn.Close()
	log.Println("Getting Channel...")

	ch, err := conn.Channel()
	if err != nil {
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return
	}

	remoteRequest := &messaging.RPCMessage{
		Type: messaging.ReadAllMessage,
		Data: []byte(""),
	}

	rpcData, err := json.Marshal(remoteRequest)
	if err != nil {
		return
	}

	corrID := randomString(32)
	err = ch.Publish(
		"",               // exchange
		service.RPCQueue, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrID,
			ReplyTo:       q.Name,
			Body:          rpcData,
		})
	if err != nil {
		return
	}
	log.Println("Listening...")
	for d := range msgs {
		rpcReply := &messaging.RPCMessage{}
		err = json.Unmarshal(d.Body, rpcReply)
		if err != nil {
			break
		}
		dogs := []*models.Dog{}
		err = json.Unmarshal(rpcReply.Data, &dogs)
		if err != nil {
			break
		}
		rpcResponse = &RPCResponse{
			Key: DogsServiceKey,
		}
		rpcResponse.Data = dogs
		break
	}

	return
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
