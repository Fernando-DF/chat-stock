package main

import (
	"fmt"
	"log"
	"net/http"

	"chat-stock/internal/bot"
	"chat-stock/internal/chat"
	"chat-stock/internal/handlers"
	"chat-stock/internal/queue"
)

func main() {
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/chat", handlers.ChatHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/ws", chat.WsHandler)

	go chat.HandleMessages()

	queue.SetupRabbitMQ()

	bot.StartBot()

	chat.ReceiveBotMessages()

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
