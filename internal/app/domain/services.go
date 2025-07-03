package domain

import (
	"lucassaraiva5/api-pay/internal/app/domain/payment"
	paypalProvider "lucassaraiva5/api-pay/internal/app/providers/paypal"
	stripeProvider "lucassaraiva5/api-pay/internal/app/providers/stripe"
)

type Services struct {
	PaymentService *payment.Service
}

func NewServices() *Services {
	paypal := paypalProvider.New()
	stripe := stripeProvider.New()
	paymentService := payment.New(paypal, stripe)
	return &Services{
		PaymentService: paymentService,
	}
}
