# Stock Chat Bot

This project implements a simple web-based chat bot that allows users to interact with a system that provides stock quotes. The bot listens for stock commands, processes them, and responds with the corresponding stock information.

## Folder Structure


## Overview

The Stock Chat Bot is a simple web application that allows users to enter stock symbols (e.g., `APPL.US`) and receive the corresponding stock information. It utilizes RabbitMQ for messaging, which allows for asynchronous communication between components. The bot listens for stock commands from the users, processes them, and responds with the requested stock quote.

### Core Components

1. **Web Server**: The main entry point of the application is the web server, which is implemented in `main.go`. This file sets up the HTTP server, routes, and serves the HTML templates (such as the login page and the chat interface).

2. **RabbitMQ Integration**: The system uses RabbitMQ for message queuing. The `rabbit_queue.go` file handles the setup and communication with RabbitMQ, including message publishing and consuming. This ensures that stock commands can be processed asynchronously.

3. **Bot Logic**: The `bot.go` file is responsible for the bot's core functionality. It listens for stock commands, retrieves the requested stock information, and sends the response back to the user.

4. **HTML Templates**: The `templates` folder contains HTML files used to render the user interface. The `login.html` template handles user authentication, while the `chat.html` template renders the chat interface where users can input stock commands.

5. **Static Assets**: The `static` folder contains CSS, JavaScript, and image files that are used for styling and interactivity in the frontend.

## Setup Instructions

### Prerequisites

Before running the application, ensure that you have the following software installed:

- Go 1.18 or later
- RabbitMQ server (for message queuing)

### Installing Dependencies

1. Clone the repository:

   ```bash
   git clone https://github.com/your_username/stock-chat-bot.git
   cd stock-chat-bot
   ```
2. Install dependencies:

   ```bash
go mod tidy
   ```


