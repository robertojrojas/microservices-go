package messaging

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type AMQManager struct {
	messageHandler *MessageHandler
	amqpURL        string
	conn           *amqp.Connection
	ch             *amqp.Channel
}

func NewAMQManager(amqpURL string, messageHandler *MessageHandler) *AMQManager {
	return &AMQManager{
		amqpURL:        amqpURL,
		messageHandler: messageHandler,
	}
}

// Connect to RabbitMQ

func (manager *AMQManager) Connect() (err error) {

	// I wonder if it would make more sense to just
	// ignore the call if connection is already open instead
	// of just closing it and opening a new one.
	if manager.conn != nil {
		manager.conn.Close() //TODO: should I handle errors here??
	}

	theConn, err := amqp.Dial(manager.amqpURL)
	if err != nil {
		return
	}
	manager.conn = theConn

	return
}

func (manager *AMQManager) Disconnect() (err error) {
	if manager.ch != nil {
		chErr := manager.ch.Close() // Ignoring this error
		log.Println(chErr)
	}

	if manager.conn != nil {
		err = manager.conn.Close()
	}

	return
}

func (manager *AMQManager) WaitForMessages() (err error) {
	msgs, err := manager.prepareToAcceptMessages()
	if err != nil {
		return
	}
	for d := range msgs {
		go manager.handleMessage(d)
	}

	return

}

func (manager *AMQManager) prepareToAcceptMessages() (msgs <-chan amqp.Delivery, err error) {

	ch, err := manager.conn.Channel()
	if err != nil {
		return
	}
	manager.ch = ch

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when usused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	if err != nil {
		return
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		return
	}

	msgs, err = ch.Consume(
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

	return

}

func (manager *AMQManager) handleMessage(delivery amqp.Delivery) {

	// Route and process incoming message
	rpcReply := manager.messageHandler.RouteMessage(delivery.Body)

	// marshal reply message
	reply, err := json.Marshal(rpcReply)
	if err != nil {
		//TODO: Need a better way to handle this
		log.Println(err)
		return
	}

	err = manager.ch.Publish(
		"",               // exchange
		delivery.ReplyTo, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: delivery.CorrelationId,
			Body:          reply,
		})

}
