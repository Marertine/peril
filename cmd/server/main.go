package main

import (
	//"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	// Create a new connection to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println("Unable to connect to RabbitMQ at localhost:5672")
	}
	defer conn.Close()

	fmt.Println("Successfully connected to RabbitMQ instance")

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	fmt.Println("Waiting for input...")
	s := <-signalChan

	// Block until a signal is received then close the channel and exit
	close(signalChan)
	fmt.Printf("Input received: %s. Shutting down\n", s)

}
