# Build stage
FROM golang:1.23 AS builder
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code and build the application
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .

# Final stage
FROM debian:bullseye-slim

# Install necessary dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# Set timezone (optional, adjust if necessary)
ENV TZ=UTC

# Copy the built binary
WORKDIR /app
COPY --from=builder /app/main .

# Ensure the binary is executable
RUN chmod +x ./main

# Run the application
CMD ["./main"]