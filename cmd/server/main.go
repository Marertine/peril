package main

import (
	//"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

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

	// Publish a message to the exchange
	err = pubsub.PublishJSON(newChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	fmt.Println("Waiting for input...")
	s := <-signalChan

	// Block until a signal is received then close the channel and exit
	close(signalChan)
	fmt.Printf("Input received: %s. Shutting down\n", s)

}
