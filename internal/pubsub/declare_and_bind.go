package pubsub

import (

	//"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	SqtDurable SimpleQueueType = iota
	SqtTransient
)

func (sqt SimpleQueueType) String() string {
	switch sqt {
	case SqtDurable:
		return "Durable"
	case SqtTransient:
		return "Transient"
	default:
		return "Unknown"
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
		return nil, amqp.Queue{}, err
	}

	// Declare a new queue
	boolDurable := (queueType == SqtDurable)
	boolAutoDelete := (queueType == SqtTransient)
	boolExclusive := (queueType == SqtTransient)
	newQueue, err := newChannel.QueueDeclare(queueName, boolDurable, boolAutoDelete, boolExclusive, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	// Bind the queue to the exchange
	err = newChannel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return newChannel, newQueue, nil

}
