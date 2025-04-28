package chat

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"chat-stock/internal/queue"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Conn     *websocket.Conn
	Username string
}

var (
	clients   = make(map[*Client]bool)
	broadcast = make(chan string)
	mu        sync.Mutex
)

// WsHandler handles the WebSocket connection upgrade and client lifecycle.
func WsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	username := cookie.Value

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		Conn:     conn,
		Username: username,
	}

	mu.Lock()
	clients[client] = true
	mu.Unlock()

	broadcast <- fmt.Sprintf("System: %s has joined the chat", username)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			mu.Lock()
			delete(clients, client)
			mu.Unlock()
			conn.Close()
			broadcast <- fmt.Sprintf("System: %s has left the chat", username)
			break
		}

		messageText := string(msg)

		if len(messageText) >= 7 && messageText[:7] == "/stock=" {
			stockCode := messageText[7:]
			log.Printf("Received stock command for: %s", stockCode)
			queue.PublishStockCommand(stockCode)
			continue
		}

		broadcast <- fmt.Sprintf("%s: %s", client.Username, messageText)
	}
}

// HandleMessages listens for incoming broadcast messages and sends to all connected clients.
func HandleMessages() {
	for {
		msg := <-broadcast

		mu.Lock()
		for client := range clients {
			err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				client.Conn.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

// ReceiveBotMessages consumes bot messages from RabbitMQ and broadcasts to chat clients.
func ReceiveBotMessages() {
	msgs, err := queue.RabbitChannel().Consume(
		"stock_messages", // queue name
		"",               // consumer
		true,             // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume stock messages: %v", err)
	}

	go func() {
		for d := range msgs {
			broadcast <- string(d.Body)
		}
	}()
}
