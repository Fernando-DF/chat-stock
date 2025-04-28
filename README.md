# Stock Chat Bot

This project implements a simple web-based chat bot that allows users to interact with a system providing real-time stock quotes.
The bot listens for stock commands, processes them asynchronously through RabbitMQ, and responds with the requested stock information.

---

## Project Structure

├── cmd
│   └── server
├── docker-compose.yml
├── Dockerfile
├── docs
│   └── go-challenge-financial-chat_5cd0c06df1e48.pdf
├── go.mod
├── go.sum
├── .env
├── internal
│   ├── bot
│   ├── chat
│   ├── handlers
│   └── queue
├── Makefile
├── README.md
├── static
└── web
    └── templates


---

## Overview

The Stock Chat Bot is a lightweight web application where users can:
- **Login** with predefined credentials.
- **Chat** in real-time using WebSocket.
- **Request stock quotes** using `/stock=CODE` (e.g., `/stock=AAPL`).
- **Receive real-time updates** from a stock bot fetching live data.

The application uses:
- **WebSocket** for real-time chat.
- **RabbitMQ** for asynchronous messaging between users and bot.
- **Go** (Golang) for backend server, bot, and message processing.

---

## Core Components

| Component     | Description |
|:--------------|:------------|
| **Web Server** | Handles HTTP routes and upgrades connections to WebSocket. |
| **Chat System** | Manages WebSocket clients and broadcasts messages. |
| **Stock Bot** | Listens to stock requests, fetches quotes, and responds back. |
| **RabbitMQ** | Manages message queues for communication between chat users and the bot. |
| **Templates** | HTML frontend for login and chat pages. |

---

## Setup Instructions

### Prerequisites

- [Go 1.20+](https://golang.org/dl/)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/)

(If you use Docker, you don't need to install RabbitMQ separately.)

---

### Local Development (without Docker)

1. Clone the repository:

```bash
git clone https://github.com/your_username/chat-stock.git
cd chat-stock
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build and run

```bash
make build
make run
```

4. Full Rebuild

```bash
make rebuild
```

### Local Development (with Docker)


4. Running:
```bash
make up
```

5. Stop and remove containers:
```bash
make down
```
