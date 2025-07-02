FROM golang:1.24-alpine

WORKDIR /app
COPY ./../.. .

RUN go mod download

WORKDIR /app/mocks/stripe
CMD ["go", "run", "stripe_mock_server.go"] 