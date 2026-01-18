package main

import (
	"fmt"
	"log"
	"strconv"

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

	myChannel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	myGameState := gamelogic.NewGameState(strUser)
	pauseHandler := handlerPause(myGameState)
	moveHandler := handlerMove(myGameState, myChannel)
	warHandler := handlerWar(myGameState, myChannel)

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

	// Subscribe to the message queue that reports wars from all players
	err = pubsub.SubscribeJSON(
		conn,                               // Connection
		routing.ExchangePerilTopic,         // Exchange
		"war",                              // QueueName
		routing.WarRecognitionsPrefix+".*", // Key
		pubsub.SqtDurable,                  // QueueType
		warHandler,                         // Handler:  func(gamelogic.HandleMove)
	)
	if err != nil {
		log.Fatalf("Unable to subscribe to queue: %s, error: %v", "war", err)
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
			if len(userInput) != 2 {
				log.Println("Spam command requires a number")
				return
			}
			loopCount, err := strconv.Atoi(userInput[1])
			if err != nil {
				log.Println("Spam command requires a number")
				return
			}
			for i := 0; i < loopCount; i++ {
				err := pubsub.PublishGameLog(myChannel, gamelogic.GetMaliciousLog(), strUser)
				if err != nil {
					log.Println("Unable to publish game log")
					return
				}
			}

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
}
