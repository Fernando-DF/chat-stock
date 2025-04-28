package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
)

var rabbitConn *amqp.Connection
var rabbitChannel *amqp.Channel
var stockQueue amqp.Queue
var stockMessagesQueue amqp.Queue

func setupRabbitMQ() {
	var err error
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	log.Printf("Connecting to RabbitMQ at %s", rabbitURL)
	rabbitConn, err = amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	rabbitChannel, err = rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	stockQueue, err = rabbitChannel.QueueDeclare(
		"stock_commands",
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare stock_commands queue: %v", err)
	}

	stockMessagesQueue, err = rabbitChannel.QueueDeclare(
		"stock_messages",
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare stock_messages queue: %v", err)
	}

	log.Println("Connected to RabbitMQ and ready.")
}

func publishStockCommand(stockCode string) {
	log.Printf("Publishing stock command for: %s", stockCode)
	err := rabbitChannel.Publish(
		"",                // exchange
		"stock_commands",  // routing key
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
