package presenter

import (
	"lucassaraiva5/api-pay/internal/app/domain/payment"
	"lucassaraiva5/api-pay/internal/app/transport/outbound"
)

func CreatePayment(payment *payment.Payment) *outbound.CreatePaymentResponse {
	return &outbound.CreatePaymentResponse{
		Id:     payment.ID,
		Status: payment.Status,
	}
}

func ReadPaymentById(payment *payment.Payment) *outbound.ReadPaymentByIdResponse {
	return &outbound.ReadPaymentByIdResponse{
		Id:     payment.ID,
		Status: payment.Status,
	}
}

func Refund(payment *payment.Payment) *outbound.RefundResponse {
	return &outbound.RefundResponse{
		Id:     payment.ID,
		Status: payment.Status,
	}
}
