# Build stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .

# Final stage
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/main .
RUN chmod +x ./main
CMD ["./main"]