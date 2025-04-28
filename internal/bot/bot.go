package bot

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"chat-stock/internal/queue"
)

// StartBot launches the stock bot that listens for stock commands from RabbitMQ.
func StartBot() {
	log.Println("Starting stock bot service...")

	msgs, err := queue.RabbitChannel().Consume(
		"stock_commands",
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
			log.Printf("Bot received stock code: %s", stockCode)

			price, err := fetchStockQuote(stockCode)
			if err != nil {
				log.Printf("Error fetching stock: %v", err)
				queue.PublishBotMessage(fmt.Sprintf("Bot: Failed to fetch stock price for %s: %v", stockCode, err))
				continue
			}

			botMsg := fmt.Sprintf("Bot: %s quote is $%s per share", strings.ToUpper(stockCode), price)
			queue.PublishBotMessage(botMsg)
		}
	}()

	log.Println("Bot started successfully and listening for commands")
}

// fetchStockQuote fetches the latest stock quote for the given stock code.
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	reader := csv.NewReader(strings.NewReader(string(body)))

	// Read CSV header
	header, err := reader.Read()
	if err != nil {
		return "", fmt.Errorf("failed to read CSV header: %v", err)
	}
	log.Printf("CSV Header: %v", header)

	// Read CSV record
	record, err := reader.Read()
	if err == io.EOF {
		return "", fmt.Errorf("no data received")
	}
	if err != nil {
		return "", fmt.Errorf("failed to read CSV: %v", err)
	}
	log.Printf("CSV Record: %v", record)

	// Find Close column
	closeIndex := -1
	for i, colName := range header {
		if strings.ToLower(colName) == "close" {
			closeIndex = i
			break
		}
	}

	if closeIndex == -1 {
		if len(record) >= 7 {
			closeIndex = 6 // Default to index 6 if no header found
		} else {
			return "", fmt.Errorf("couldn't determine close price column and record doesn't have enough columns")
		}
	}

	if len(record) <= closeIndex {
		return "", fmt.Errorf("CSV data does not contain enough columns")
	}

	price := record[closeIndex]
	if price == "N/D" {
		return "", fmt.Errorf("invalid stock code or no data available")
	}

	return price, nil
}
