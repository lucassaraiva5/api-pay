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
# RUN go test ./... -v