package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func SendMessage() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Create a Queue to send the message to.
	q, err := ch.QueueDeclare(
		"golang-queue", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Let's catch the message from the terminal
	var mPayload string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Please enter message: ")
	if scanner.Scan() {
		mPayload = scanner.Text()
	}

	// Set the payload for the message
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(mPayload),
		},
	)

	// If there is an error publishing the message, a log will be displayed on terminal
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Congrats, sending message: %s", mPayload)
}
