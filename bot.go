package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"github.com/streadway/amqp"
)

func startBot() {
	log.Println("Starting stock bot service...")

	msgs, err := rabbitChannel.Consume(
		stockQueue.Name, // "stock_commands" queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer for stock commands: %v", err)
	}

	go func() {
		for d := range msgs {
			stockCode := string(d.Body)
			log.Printf("Bot received stock code: %s", stockCode)

			price, err := fetchStockQuote(stockCode)
			if err != nil {
				log.Printf("Error fetching stock: %v", err)
				publishBotMessage(fmt.Sprintf("Bot: Failed to fetch stock price for %s: %v", stockCode, err))
				continue
			}

			botMsg := fmt.Sprintf("Bot: %s quote is $%s per share", strings.ToUpper(stockCode), price)
			publishBotMessage(botMsg)
		}
	}()

	log.Println("Bot started successfully and listening for commands")
}

func fetchStockQuote(stockCode string) (string, error) {
	url := fmt.Sprintf("https://stooq.com/q/l/?s=%s.us&f=sd2t2ohlcv&h&e=csv", stockCode)
	log.Printf("Fetching stock data from URL: %s", url)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch quote: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned non-OK status: %d %s", resp.StatusCode, resp.Status)
	}

	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Uncomment for debugging if needed
	// log.Printf("API Response body: %s", string(body))

	// Parse as CSV
	reader := csv.NewReader(strings.NewReader(string(body)))

	// Read header
	header, err := reader.Read()
	if err != nil {
		return "", fmt.Errorf("failed to read CSV header: %v", err)
	}
	log.Printf("CSV Header: %v", header)

	// Read data record
	record, err := reader.Read()
	if err == io.EOF {
		return "", fmt.Errorf("no data received")
	}
	if err != nil {
		return "", fmt.Errorf("failed to read CSV: %v", err)
	}

	log.Printf("CSV Record: %v", record)

	// Find the Close price column
	closeIndex := -1
	for i, colName := range header {
		if strings.ToLower(colName) == "close" {
			closeIndex = i
			break
		}
	}

	// If we couldn't find the Close column, use a default index (6 is typical)
	if closeIndex == -1 {
		if len(record) >= 7 {
			closeIndex = 6
		} else {
			return "", fmt.Errorf("couldn't determine close price column and record doesn't have enough columns")
		}
	}

	// Check if we have enough columns in the record
	if len(record) <= closeIndex {
		return "", fmt.Errorf("CSV data does not contain enough columns")
	}

	price := record[closeIndex]
	if price == "N/D" {
		return "", fmt.Errorf("invalid stock code or no data available")
	}

	return price, nil
}

func publishBotMessage(message string) {
	log.Printf("Publishing bot message: %s", message)
	err := rabbitChannel.Publish(
		"",              // exchange
		"stock_messages", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Printf("Failed to publish bot message: %v", err)
	}
}
