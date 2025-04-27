package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/streadway/amqp"
)

func startBot() {
	msgs, err := rabbitChannel.Consume(
		stockQueue.Name, // "stock_commands" queue
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer for stock commands: %v", err)
	}

	go func() {
		for d := range msgs {
			stockCode := string(d.Body)
			fmt.Printf("Bot received stock code: %s\n", stockCode)

			price, err := fetchStockQuote(stockCode)
			if err != nil {
				publishBotMessage(fmt.Sprintf("Bot: Failed to fetch stock price for %s", stockCode))
				continue
			}

			botMsg := fmt.Sprintf("Bot: %s quote is $%s per share", strings.ToUpper(stockCode), price)
			publishBotMessage(botMsg)
		}
	}()
}

func fetchStockQuote(stockCode string) (string, error) {
	url := fmt.Sprintf("https://stooq.com/q/l/?s=%s.us&f=sd2t2ohlcv&h&e=csv", stockCode)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch quote: %v", err)
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	_, _ = reader.Read() // skip header

	record, err := reader.Read()
	if err == io.EOF {
		return "", fmt.Errorf("no data received")
	}
	if err != nil {
		return "", fmt.Errorf("failed to read CSV: %v", err)
	}

	// CSV format: Symbol, Date, Time, Open, High, Low, Close, Volume
	price := record[6]
	if price == "N/D" {
		return "", fmt.Errorf("invalid stock code")
	}

	return price, nil
}

func publishBotMessage(message string) {
	err := rabbitChannel.Publish(
		"", // exchange
		"stock_messages",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Printf("Failed to publish bot message: %v", err)
	}
}
