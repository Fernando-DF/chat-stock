package queue

import (
	"log"
	"os"
	"time"

	"github.com/streadway/amqp"
)

var (
	rabbitConn         *amqp.Connection
	rabbitChannel      *amqp.Channel
	stockQueue         amqp.Queue
	stockMessagesQueue amqp.Queue
)

// SetupRabbitMQ initializes the connection to RabbitMQ and declares queues.
func SetupRabbitMQ() {
	var err error
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	// Retry logic
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to connect to RabbitMQ at %s (attempt %d/%d)", rabbitURL, i+1, maxRetries)
		rabbitConn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ after %d attempts: %v", maxRetries, err)
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

// PublishStockCommand sends a stock code to the stock_commands queue.
func PublishStockCommand(stockCode string) {
	log.Printf("Publishing stock command for: %s", stockCode)
	err := rabbitChannel.Publish(
		"",
		"stock_commands",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(stockCode),
		},
	)
	if err != nil {
		log.Printf("Failed to publish stock command: %v", err)
	}
}

// PublishBotMessage sends a bot message to the stock_messages queue.
func PublishBotMessage(message string) {
	log.Printf("Publishing bot message: %s", message)
	err := rabbitChannel.Publish(
		"",
		"stock_messages",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Printf("Failed to publish bot message: %v", err)
	}
}

// RabbitChannel provides access to the underlying RabbitMQ channel (for consumers).
func RabbitChannel() *amqp.Channel {
	return rabbitChannel
}
