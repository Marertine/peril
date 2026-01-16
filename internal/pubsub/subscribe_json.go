package pubsub

import (
	//"context"
	"encoding/json"
	"fmt"
	"log"

	//"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	//"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

const (
	Ack Acktype = iota
	NackDiscard
	NackRequeue
)

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) Acktype) error {
	// Ensure that the queue exists and is bound to the exchange
	chanDelivery, queueDelivery, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	msgs, err := chanDelivery.Consume(queueDelivery.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		defer chanDelivery.Close()
		for msg := range msgs {
			var v T
			err := json.Unmarshal(msg.Body, &v)
			if err != nil {
				log.Printf("failed to unmarshal message %s: %v", msg.MessageId, err)
				continue
			}

			handler(v)
			err = msg.Ack(false)
			if err != nil {
				log.Printf("Unable to acknowledge messageID: %s: %v", msg.MessageId, err)
				continue
			}

			switch handler(v) {
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
