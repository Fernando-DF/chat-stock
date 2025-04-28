FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for go get
RUN apk add --no-cache git

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o chat-app .

# Final stage
FROM alpine:3.18

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/chat-app .

# Copy templates and static folders if they exist
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Set environment variables
ENV RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./chat-app"]
