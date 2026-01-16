package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"
	fmt.Println("Starting Peril server...")

	// Create a new connection to RabbitMQ
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("Unable to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	fmt.Println("Peril game server connected to RabbitMQ")

	// Create a messages channel
	msgChannel, err := conn.Channel()
	if err != nil {
		log.Fatal("Unable to create messages queue")
	}

	// Create a logs queue
	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.SqtDurable)
	if err != nil {
		log.Fatalf("Unable to declare and bind: %v", err)
	}
	fmt.Println("Logs queue created.")

	// Print the server commands to the console
	gamelogic.PrintServerHelp()

	// Wait for user input
	for {
		userInput := gamelogic.GetInput()

		switch userInput[0] {
		case "help":
			gamelogic.PrintServerHelp()

		case "pause":
			helperServerPause(msgChannel, true)

		case "quit":
			fmt.Println("Shutting down at user request.")
			return

		case "resume":
			helperServerPause(msgChannel, false)

		default:
			log.Printf("Unrecognised command: %s\n", userInput[0])
		}
	}

}

func helperServerPause(channel *amqp.Channel, pause bool) error {
	// Publish a message to the exchange
	err := pubsub.PublishJSON(
		channel,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{
			IsPaused: pause,
		},
	)
	if err != nil {
		return err
	}

	if pause {
		fmt.Println("Pause message sent!")
	} else {
		fmt.Println("Resume message sent!")
	}

	return nil
}
