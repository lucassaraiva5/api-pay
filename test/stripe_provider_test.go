package test

import (
	stripeProvider "lucassaraiva5/api-pay/internal/app/providers/stripe"
	"testing"

	"lucassaraiva5/api-pay/internal/app/domain/stripe"
	"lucassaraiva5/api-pay/internal/app/transport/outbound"

	"github.com/google/uuid"
)

func TestStripe_CreatePayment_Success(t *testing.T) {
	provider := stripeProvider.New()

	request := &stripe.PaymentRequest{
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

	payment := response.(*outbound.StripePaymentResponse)
	if payment.Status != "paid" {
		t.Fatalf("expected status to be paid, got %s", payment.Status)
	}

	if _, err := uuid.Parse(payment.ID); err != nil {
		t.Fatalf("expected a valid UUID, got %s", payment.ID)
	}
}

func TestStripe_CreatePayment_InvalidRequest(t *testing.T) {
	provider := stripeProvider.New()

	request := struct {
		InvalidField string
	}{
		InvalidField: "invalid",
	}

	_, err := provider.CreatePayment(request)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if err.Error() != "invalid request type for Stripe" {
		t.Fatalf("expected error 'invalid request type for Stripe', got %v", err)
	}
}

func TestStripe_Refund_Success(t *testing.T) {
	provider := stripeProvider.New()

	request := &stripe.PaymentRequest{
		Amount:   100.0,
		Currency: "USD",
	}
	response, _ := provider.CreatePayment(request)
	payment := response.(*outbound.StripePaymentResponse)

	refundResponse, err := provider.Refund(payment.ID, 50.0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	refundedPayment := refundResponse.(*outbound.StripePaymentResponse)
	if refundedPayment.Status != "voided" {
		t.Fatalf("expected status to be voided, got %s", refundedPayment.Status)
	}

	if refundedPayment.Amount != 50.0 {
		t.Fatalf("expected remaining amount to be 50.0, got %f", refundedPayment.Amount)
	}
}

func TestStripe_Refund_PaymentNotFound(t *testing.T) {
	provider := stripeProvider.New()

	_, err := provider.Refund("invalid-id", 50.0)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestStripe_GetPayment_Success(t *testing.T) {
	provider := stripeProvider.New()

	request := &stripe.PaymentRequest{
		Amount:   100.0,
		Currency: "USD",
	}
	response, _ := provider.CreatePayment(request)
	payment := response.(*outbound.StripePaymentResponse)

	getResponse, err := provider.GetPayment(payment.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrievedPayment := getResponse.(*outbound.StripePaymentResponse)
	if retrievedPayment.ID != payment.ID {
		t.Fatalf("expected payment ID to match, got %s", retrievedPayment.ID)
	}
}

func TestStripe_GetPayment_PaymentNotFound(t *testing.T) {
	provider := stripeProvider.New()

	_, err := provider.GetPayment("invalid-id")
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}
