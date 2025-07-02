FROM golang:1.24-alpine

WORKDIR /app
COPY ./../.. .

RUN go mod download

WORKDIR /app/mocks/paypal
CMD ["go", "run", "paypal_mock_server.go"] 