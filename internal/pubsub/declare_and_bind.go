package pubsub

import (
	//"context"
	//"encoding/json"

	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	sqtDurable SimpleQueueType = iota
	sqtTransient
)

func (sqt SimpleQueueType) String() string {
	switch sqt {
	case sqtDurable:
		return "Durable"
	case sqtTransient:
		return "Transient"
	}
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {

	// Create a new channel
	newChannel, err := conn.Channel()
	if err != nil {
		log.Fatal("Unable to create new channel")
	}

	// Declare a new queue
	boolDurable := (queueType == sqtDurable)
	boolAutoDelete := (queueType == sqtTransient)
	boolExclusive := (queueType == sqtTransient)
	newQueue, err := newChannel.QueueDeclare("name", boolDurable, boolAutoDelete, boolExclusive, false, nil)
	if err != nil {
		log.Fatal("Unable to create new queue")
	}

	// Bind the queue to the exchange
	err = newChannel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		log.Fatal("Unable to create bind new queue")
	}

	return newChannel, newQueue, nil

}
