package main

import (
	//"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	//"os"
	//"os/signal"

	//"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"

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

	strUser, err := ClientWelcome()

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, strUser)
	myChannel, myQueue, err := DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.sqtTransient)
	if err != nil {
		log.Fatalf("Unable to declare and bind: %v", err)
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	//fmt.Println("Waiting for input...")
	s := <-signalChan

}
