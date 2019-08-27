package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func GetChannel(url string) (*amqp.Connection, * amqp.Channel) {
	conn, err := amqp.Dial(url)
	failOnError(err, "failed to establish connection to message broker")
	ch, err := conn.Channel()
	failOnError(err, "failed to get channel for connection")

	return conn, ch
}

func getQueue(name string, ch *amqp.Channel) *amqp.Queue {
	q, err := ch.QueueDeclare(name, false, false, false, false, nil)
	failOnError(err, "failed to declare queue")
}

func failOnError(err error, msg string){
	if err != nil {
		log.Fatalf("%s: %s",msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}