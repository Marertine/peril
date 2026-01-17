package pubsub

import (
	//"bytes"
	//"context"
	//"encoding/gob"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGameLog(channel *amqp.Channel, message, username string) error {

	key := routing.GameLogSlug + "." + username
	exchange := routing.ExchangePerilTopic

	data := routing.GameLog{
		CurrentTime: time.Now().UTC(),
		Message:     message,
		Username:    username,
	}

	return PublishGob(channel, exchange, key, data)

}
