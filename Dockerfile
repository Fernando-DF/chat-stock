# Dockerfile
FROM golang:1.24-alpine

# Install necessary tools
RUN apk --no-cache add git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app
COPY . .

# Build the Go app
RUN go build -o app .

# Expose the port your app uses
EXPOSE 8080

# Command to run the app
CMD ["./app"]
