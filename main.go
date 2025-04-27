package main

import (
	"fmt"
	"html/template"
	"net/http"
	"github.com/gorilla/websocket"
	"sync"
	"log"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn     *websocket.Conn
	username string
}

var (
	tpl        = template.Must(template.ParseGlob("templates/*.html"))
	users      = map[string]string{
		"marcelo": "secret123",
		"admin":   "admin123",
	} // hardcoded user for now

	sessionMap = map[string]string{} // sessionID -> username
	clients   = make(map[*Client]bool)
	broadcast = make(chan string)
	mu        sync.Mutex
)

func generateSessionID(username string) string {
	return "session_" + username
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	if users[username] == password {
		sessionID := generateSessionID(username)
		sessionMap[sessionID] = username

		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: sessionID,
			Path:  "/",
		})
		http.Redirect(w, r, "/chat", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "login.html", "Invalid credentials")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == nil {
		delete(sessionMap, cookie.Value)

		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || sessionMap[cookie.Value] == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tpl.ExecuteTemplate(w, "chat.html", sessionMap[cookie.Value])
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || sessionMap[cookie.Value] == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	username := sessionMap[cookie.Value]

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	client := &Client{conn: conn, username: username}

	mu.Lock()
	clients[client] = true
	mu.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			mu.Lock()
			delete(clients, client)
			mu.Unlock()
			conn.Close()
			break
		}

		messageText := string(msg)

	if len(messageText) >= 7 && messageText[:7] == "/stock=" {
		stockCode := messageText[7:]
		fmt.Println("Received stock command for:", stockCode)

		publishStockCommand(stockCode)

		continue
	}

		broadcast <- fmt.Sprintf("%s: %s", client.username, string(msg))
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		mu.Lock()
		for client := range clients {
			err := client.conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				client.conn.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func listenBotMessages() {
	msgs, err := rabbitChannel.Consume(
		"stock_messages",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer for bot messages: %v", err)
	}

	go func() {
		for d := range msgs {
			botMessage := string(d.Body)
			broadcast <- botMessage
		}
	}()
}

func consumeBotMessages() {
	msgs, err := rabbitChannel.Consume(
		"stock_messages", // queue
		"",             // consumer
		true,           // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		log.Fatalf("Failed to consume stock messages: %v", err)
	}

	go func() {
		for d := range msgs {
			// d.Body contains the message from the bot
			mu.Lock()
			for client := range clients {
				err := client.conn.WriteMessage(websocket.TextMessage, d.Body)
				if err != nil {
					client.conn.Close()
					delete(clients, client)
				}
			}
			mu.Unlock()
		}
	}()
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/chat", chatHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/ws", wsHandler)
	go handleMessages()

	setupRabbitMQ()
	startBot()
	listenBotMessages()
	consumeBotMessages()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
