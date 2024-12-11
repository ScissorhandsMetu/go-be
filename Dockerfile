# Use the official Go image to create a build environment
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal base image for the final container
FROM alpine:latest

# Install necessary libraries
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the built Go application from the builder stage
COPY --from=builder /app/main .

# Expose port 8000
EXPOSE 8000

# Command to run the application
CMD ["./main"]