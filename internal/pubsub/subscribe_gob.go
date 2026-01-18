package pubsub

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGob[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) Acktype) error {
	// Ensure that the queue exists and is bound to the exchange
	chanDelivery, queueDelivery, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	chanDelivery.Qos(10, 0, false)

	msgs, err := chanDelivery.Consume(queueDelivery.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		defer chanDelivery.Close()
		for msg := range msgs {
			network := bytes.NewBuffer(msg.Body)
			dec := gob.NewDecoder(network)
			var v T

			err := dec.Decode(&v)
			if err != nil {
				log.Printf("failed to decode message %s: %v", msg.MessageId, err)
				continue
			}

			ackType := handler(v)
			switch ackType {
			case Ack:
				msg.Ack(false)
				fmt.Println("Ack")
			case NackDiscard:
				msg.Nack(false, false)
				fmt.Println("NackDiscard")
			case NackRequeue:
				msg.Nack(false, true)
				fmt.Println("NackRequeue")
			}

		}
	}()

	return nil

}
