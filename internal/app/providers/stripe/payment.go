package stripeProvider

import (
	"errors"
	"lucassaraiva5/api-pay/internal/app/domain/stripe"
	"lucassaraiva5/api-pay/internal/app/transport/outbound"

	"github.com/google/uuid"
)

type Provider struct {
	Payments map[string]outbound.StripePaymentResponse
}

func New() *Provider {
	return &Provider{
		Payments: make(map[string]outbound.StripePaymentResponse),
	}
}

func (p *Provider) CreatePayment(req interface{}) (interface{}, error) {
	request, ok := req.(*stripe.PaymentRequest)
	if !ok {
		return nil, errors.New("invalid request type for Stripe")
	}

	id := uuid.New().String()
	payment := outbound.StripePaymentResponse{
		ID:          id,
		Status:      "paid",
		Amount:      request.Amount,
		Description: request.Description,
		PaymentType: request.PaymentType,
		Currency:    request.Currency,
		Card:        request.Card,
	}
	p.Payments[id] = payment
	return &payment, nil
}

func (p *Provider) Refund(paymentID string, amount float64) (interface{}, error) {
	payment, exists := p.Payments[paymentID]
	if !exists {
		return nil, errors.New("payment not found")
	}

	payment.Status = "voided"
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
