# Use the official Go image for building the application
FROM golang:1.24.3-alpine AS builder

# Install additional dependencies
RUN apk add --no-cache git bash

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the project files
COPY . .

# Run unit tests to validate the code
RUN go test ./... -v

# Build the application binary from the main.go file located in the cmd directory
RUN go build -o main ./cmd

# Use a lightweight Alpine image for the final stage
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk add --no-cache bash

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main /app/main

# Set the working directory inside the container
WORKDIR /app

# Expose the application port
EXPOSE 8088

# Define the default command to run the application
CMD ["./main"]