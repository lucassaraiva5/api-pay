package outbound

import (
	"lucassaraiva5/api-pay/internal/app/domain/stripe"
)

type StripePaymentResponse struct {
	ID          string      `json:"id"`
	Status      string      `json:"status"`
	Amount      float64     `json:"amount"`
	Description string      `json:"description"`
	PaymentType string      `json:"payment_type"`
	CreatedAt   string      `json:"created_at"`
	Currency    string      `json:"currency"`
	Card        stripe.Card `json:"card"`
}
