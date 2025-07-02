package paypalProvider

import (
	"errors"
	"lucassaraiva5/api-pay/internal/app/domain/paypal"
	"lucassaraiva5/api-pay/internal/app/transport/outbound"
	"time"

	"github.com/google/uuid"
)

type Provider struct {
	Payments map[string]outbound.PaypalPaymentResponse
}

func New() *Provider {
	return &Provider{
		Payments: make(map[string]outbound.PaypalPaymentResponse),
	}
}

func (p *Provider) CreatePayment(req interface{}) (interface{}, error) {
	request, ok := req.(*paypal.PaymentRequest)
	if !ok {
		return nil, errors.New("invalid request type for PayPal")
	}

	id := uuid.New().String()
	payment := outbound.PaypalPaymentResponse{
		ID:          id,
		Status:      "authorized",
		Amount:      request.Amount,
		Currency:    request.Currency,
		Description: request.Description,
		CreatedAt:   time.DateTime,
		Method:      request.PaymentMethod,
	}
	p.Payments[id] = payment
	return &payment, nil
}

func (p *Provider) Refund(paymentID string, amount float64) (interface{}, error) {
	payment, exists := p.Payments[paymentID]
	if !exists {
		return nil, errors.New("payment not found")
	}

	payment.Status = "refunded"
	payment.Amount -= amount
	p.Payments[paymentID] = payment
	return &payment, nil
}

func (p *Provider) GetPayment(paymentID string) (interface{}, error) {
	payment, exists := p.Payments[paymentID]
	if !exists {
		return nil, errors.New("payment not found")
	}
	return &payment, nil
}
