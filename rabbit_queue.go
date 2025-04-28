package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"time"
)

var rabbitConn *amqp.Connection
var rabbitChannel *amqp.Channel
var stockQueue amqp.Queue
var stockMessagesQueue amqp.Queue

func setupRabbitMQ() {
	var err error
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	
	// Add retry logic for RabbitMQ connection
	var conn *amqp.Connection
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to connect to RabbitMQ at %s (attempt %d/%d)", rabbitURL, i+1, maxRetries)
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}
	
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ after %d attempts: %v", maxRetries, err)
	}
	
	rabbitConn = conn
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
