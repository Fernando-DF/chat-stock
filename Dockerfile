FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.18

WORKDIR /app

# Copy only the binary and needed files
COPY --from=builder /app/server .
COPY --from=builder /app/web/templates ./web/templates

# Set environment variables
ENV RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/

EXPOSE 8080

CMD ["./server"]

