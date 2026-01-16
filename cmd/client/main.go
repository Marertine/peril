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

	//queueName := fmt.Sprintf("%s.%s", routing.PauseKey, strUser)

	myGameState := gamelogic.NewGameState(strUser)
	pauseHandler := handlerPause(myGameState)
	moveHandler := handlerMove(myGameState)

	// Subscribe to users own message queue
	err = pubsub.SubscribeJSON(
		conn,                         // Connection
		routing.ExchangePerilDirect,  // Exchange
		routing.PauseKey+"."+strUser, // QueueName
		routing.PauseKey,             // Key
		pubsub.SqtTransient,          // QueueType
		pauseHandler,                 // Handler: func(routing.PlayingState)
	)
	if err != nil {
		log.Fatalf("Unable to subscribe to queue: %s, error: %v", routing.PauseKey+"."+strUser, err)
	}

	// Subscribe to the message queue that reports moves from all players
	err = pubsub.SubscribeJSON(
		conn,                                // Connection
		routing.ExchangePerilTopic,          // Exchange
		routing.ArmyMovesPrefix+"."+strUser, // QueueName
		routing.ArmyMovesPrefix+".*",        // Key
		pubsub.SqtTransient,                 // QueueType
		moveHandler,                         // Handler:  func(gamelogic.HandleMove)
	)
	if err != nil {
		log.Fatalf("Unable to subscribe to queue: %s, error: %v", routing.ArmyMovesPrefix+".*", err)
	}

	myChannel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}

	// Wait for user input
	for {
		userInput := gamelogic.GetInput()

		switch userInput[0] {
		case "help":
			gamelogic.PrintClientHelp()

		case "move":
			myArmyMove, err := myGameState.CommandMove(userInput)
			if err != nil {
				fmt.Println(err)
			}

			err = pubsub.PublishJSON(myChannel, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+"."+strUser, myArmyMove)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Move published successfully!")
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
