APP_NAME=server

# Try to detect docker-compose version (docker compose vs docker-compose)
DOCKER_COMPOSE=$(shell command -v docker-compose >/dev/null 2>&1 && echo docker-compose || echo docker compose)

include .env
export $(shell sed 's/=.*//' .env)

build:
	go build -o $(APP_NAME) ./cmd/server

run: build
	./$(APP_NAME)

up:
	$(DOCKER_COMPOSE) --env-file .env up --build

down:
	$(DOCKER_COMPOSE) down

clean:
	rm -f $(APP_NAME)

rebuild:
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) build --no-cache
	$(DOCKER_COMPOSE) --env-file .env up --build
