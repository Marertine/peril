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

	// Create a new channel
	newChannel, err := conn.Channel()
	if err != nil {
		log.Fatal("Unable to create new channel")
	}

	// Print the server commands to the console
	gamelogic.PrintServerHelp()

	// Wait for user input
	for {
		userInput := gamelogic.GetInput()

		switch userInput[0] {
		case "help":
			gamelogic.PrintServerHelp()

		case "pause":
			helperServerPause(newChannel, true)

		case "quit":
			fmt.Println("Shutting down at user request.")
			return

		case "resume":
			helperServerPause(newChannel, false)

		default:
			log.Printf("Unrecognised command: %s\n", userInput[0])
		}
	}

	// It appears the tasks beyond this point are no longer intended to be used, but the lesson did not stipulate as such
	/*
		// wait for ctrl+c
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		fmt.Println("Waiting for input...")
		s := <-signalChan

		// Block until a signal is received then close the channel and exit
		close(signalChan)
		fmt.Printf("Input received: %s. Shutting down\n", s)
	*/

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
