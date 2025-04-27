package main

import (
	"github.com/streadway/amqp"
	"log"
)

var rabbitConn *amqp.Connection
var rabbitChannel *amqp.Channel
var stockQueue amqp.Queue
var stockMessagesQueue amqp.Queue

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
		"stock_commands",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare stock_commands queue: %v", err)
	}

	stockMessagesQueue, err = rabbitChannel.QueueDeclare(
		"stock_messages",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare stock_messages queue: %v", err)
	}

	log.Println("Connected to RabbitMQ and ready.")
}

func publishStockCommand(stockCode string) {
	err := rabbitChannel.Publish(
		"",
		"stock_commands",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(stockCode),
		})
	if err != nil {
		log.Printf("Failed to publish stock command: %v", err)
	}
}
