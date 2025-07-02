package mapper

import (
	"lucassaraiva5/api-pay/internal/app/domain/payment"
	"lucassaraiva5/api-pay/internal/app/transport/inbound"
)

func PaymentFromCreatePaymentRequest(request *inbound.CreatePaymentRequest) *payment.Payment {
	return &payment.Payment{
		Amount:      request.Amount,
		Currency:    request.Currency,
		Description: request.Description,
		Method: payment.Method{
			Type: request.Method.Type,
			Card: payment.Card{
				Number:            request.Method.Card.Number,
				Holder:            request.Method.Card.Holder,
				CVV:               request.Method.Card.CVV,
				Expiration:        request.Method.Card.Expiration,
				InstallmentNumber: request.Method.Card.InstallmentNumber,
			},
		},
	}
}
