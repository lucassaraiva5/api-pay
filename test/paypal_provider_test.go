package test

import (
	paypalProvider "lucassaraiva5/api-pay/internal/app/providers/paypal"
	"testing"

	"lucassaraiva5/api-pay/internal/app/domain/paypal"
	"lucassaraiva5/api-pay/internal/app/transport/outbound"

	"github.com/google/uuid"
)

func TestPayPal_CreatePayment_Success(t *testing.T) {
	provider := paypalProvider.New()

	request := &paypal.PaymentRequest{
		Amount:   100.0,
		Currency: "USD",
	}

	response, err := provider.CreatePayment(request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response == nil {
		t.Fatalf("expected a valid response, got nil")
	}

	payment := response.(*outbound.PaypalPaymentResponse)
	if payment.Status != "authorized" {
		t.Fatalf("expected status to be authorized, got %s", payment.Status)
	}

	if _, err := uuid.Parse(payment.ID); err != nil {
		t.Fatalf("expected a valid UUID, got %s", payment.ID)
	}
}

func TestPayPal_CreatePayment_InvalidRequest(t *testing.T) {
	provider := paypalProvider.New()

	request := struct {
		InvalidField string
	}{
		InvalidField: "invalid",
	}

	_, err := provider.CreatePayment(request)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if err.Error() != "invalid request type for PayPal" {
		t.Fatalf("expected error 'invalid request type for PayPal', got %v", err)
	}
}

func TestPayPal_Refund_Success(t *testing.T) {
	provider := paypalProvider.New()

	request := &paypal.PaymentRequest{
		Amount:   100.0,
		Currency: "USD",
	}
	response, _ := provider.CreatePayment(request)
	payment := response.(*outbound.PaypalPaymentResponse)

	refundResponse, err := provider.Refund(payment.ID, 50.0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	refundedPayment := refundResponse.(*outbound.PaypalPaymentResponse)
	if refundedPayment.Status != "refunded" {
		t.Fatalf("expected status to be refunded, got %s", refundedPayment.Status)
	}

	if refundedPayment.Amount != 50.0 {
		t.Fatalf("expected remaining amount to be 50.0, got %f", refundedPayment.Amount)
	}
}

func TestPayPal_Refund_PaymentNotFound(t *testing.T) {
	provider := paypalProvider.New()

	_, err := provider.Refund("invalid-id", 50.0)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestPayPal_GetPayment_Success(t *testing.T) {
	provider := paypalProvider.New()

	request := &paypal.PaymentRequest{
		Amount:   100.0,
		Currency: "USD",
	}
	response, _ := provider.CreatePayment(request)
	payment := response.(*outbound.PaypalPaymentResponse)

	getResponse, err := provider.GetPayment(payment.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrievedPayment := getResponse.(*outbound.PaypalPaymentResponse)
	if retrievedPayment.ID != payment.ID {
		t.Fatalf("expected payment ID to match, got %s", retrievedPayment.ID)
	}
}

func TestPayPal_GetPayment_PaymentNotFound(t *testing.T) {
	provider := paypalProvider.New()

	_, err := provider.GetPayment("invalid-id")
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}
