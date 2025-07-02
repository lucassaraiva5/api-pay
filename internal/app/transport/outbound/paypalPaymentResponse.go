package outbound

import (
	"lucassaraiva5/api-pay/internal/app/domain/paypal"
)

type PaypalPaymentResponse struct {
	ID          string               `json:"id"`
	Status      string               `json:"status"`
	Amount      float64              `json:"amount"`
	Currency    string               `json:"currency"`
	Description string               `json:"description"`
	CreatedAt   string               `json:"created_at"`
	Method      paypal.PaymentMethod `json:"method"`
}
