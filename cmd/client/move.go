package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, channel *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(move gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(move, channel)

		switch moveOutcome {
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard

		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack

		case gamelogic.MoveOutcomeMakeWar:
			key := routing.WarRecognitionsPrefix + "." + gs.GetUsername()
			publishData := gamelogic.RecognitionOfWar{
				Attacker: move.Player,
				Defender: gs.GetPlayerSnap(),
			}
			err := pubsub.PublishJSON(channel, routing.ExchangePerilTopic, key, publishData)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Move published successfully!")
			}
			return pubsub.NackRequeue
		}
		fmt.Println("error: unknown move outcome")
		return pubsub.NackDiscard
	}
}

func handlerWar(gs *gamelogic.GameState, channel *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(recognition gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")
		//Call the gamestate's HandleWar method with the message's body (recognition).
		outcome, _, _ := gs.HandleWar(recognition)

		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue

		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard

		case gamelogic.WarOutcomeOpponentWon:
			return pubsub.Ack

		case gamelogic.WarOutcomeYouWon:
			return pubsub.Ack

		case gamelogic.WarOutcomeDraw:
			return pubsub.Ack

		default:
			fmt.Println("Unrecognised war outcome!")
			return pubsub.NackDiscard
		}

	}
}
