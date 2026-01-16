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

	fmt.Println("Starting Peril client...")

	// Create a new connection to RabbitMQ
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("Unable to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	strUser, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("No valid user provided: %v", err)
	}

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, strUser)
	/*_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.SqtTransient)
	if err != nil {
		log.Fatalf("Unable to declare and bind: %v", err)
	}*/

	myGameState := gamelogic.NewGameState(strUser)
	pauseHandler := handlerPause(myGameState)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		queueName,
		routing.PauseKey,
		pubsub.SqtTransient,
		pauseHandler, // this is func(routing.PlayingState)
	)

	// Wait for user input
	for {
		userInput := gamelogic.GetInput()

		switch userInput[0] {
		case "help":
			gamelogic.PrintClientHelp()

		case "move":
			_, err := myGameState.CommandMove(userInput)
			if err != nil {
				fmt.Println(err)
			}

		case "quit":
			gamelogic.PrintQuit()
			return

		case "spam":
			log.Println("Spamming not allowed yet!")

		case "spawn":
			err := myGameState.CommandSpawn(userInput)
			if err != nil {
				fmt.Println(err)
			}

		case "status":
			myGameState.CommandStatus()

		default:
			log.Printf("Unrecognised command: %s\n", userInput[0])
		}
	}

	/*
		// wait for ctrl+c
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		_ = <-signalChan
	*/
}
