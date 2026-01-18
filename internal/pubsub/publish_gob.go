package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var network bytes.Buffer        // Buffer that will hold the encoded data
	enc := gob.NewEncoder(&network) // Encoder that explains how to write gob-encoded data into the buffer
	err := enc.Encode(val)          // Perform the encoding of val and actually write the bytes into network

	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/gob",
			Body:        network.Bytes(),
		},
	)
	return err
}
