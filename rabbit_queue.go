package main

import (
	"github.com/streadway/amqp"
	"log"
)

var rabbitConn *amqp.Connection
var rabbitChannel *amqp.Channel
var stockQueue amqp.Queue

func setupRabbitMQ() {
	var err error
	rabbitConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	rabbitChannel, err = rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	stockQueue, err = rabbitChannel.QueueDeclare(
		"stock_commands", // Queue name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	log.Println("Connected to RabbitMQ and ready.")
}

func publishStockCommand(stockCode string) {
	err := rabbitChannel.Publish(
		"",                // exchange
		stockQueue.Name,   // routing key (queue name)
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(stockCode),
		})
	if err != nil {
		log.Printf("Failed to publish stock command: %v", err)
	}
}
