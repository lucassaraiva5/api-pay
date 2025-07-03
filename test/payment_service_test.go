package test

import (
	"context"
	"lucassaraiva5/api-pay/internal/app/domain/payment"
	paypalProvider "lucassaraiva5/api-pay/internal/app/providers/paypal"
	stripeProvider "lucassaraiva5/api-pay/internal/app/providers/stripe"
	"testing"

	"github.com/google/uuid"
)

func TestProcessPayment_SuccessWithPrimaryProvider(t *testing.T) {
	primary := paypalProvider.New()
	secondary := stripeProvider.New()

	service := &payment.Service{
		PrimaryProvider:   primary,
		SecondaryProvider: secondary,
	}

	paymentRequest := &payment.Payment{
		Amount:   100.0,
		Currency: "USD",
		Method: payment.Method{
			Type: "card",
			Card: payment.Card{
				Number:            "4111111111111111",
				Holder:            "John Doe",
				CVV:               "123",
				Expiration:        "12/2025",
				InstallmentNumber: 1,
			},
		},
	}

	result, err := service.ProcessPayment(context.Background(), paymentRequest)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil || result.Status != "authorized" {
		t.Fatalf("expected payment to be authorized, got %v", result)
	}

	if _, err := uuid.Parse(result.ID); err != nil {
		t.Fatalf("expected a valid UUID, got %s", result.ID)
	}
}

func TestRefundPayment_SuccessWithPrimaryProvider(t *testing.T) {
	primary := paypalProvider.New()
	secondary := stripeProvider.New()

	service := &payment.Service{
		PrimaryProvider:   primary,
		SecondaryProvider: secondary,
	}

	paymentRequest := &payment.Payment{
		Amount:   100.0,
		Currency: "USD",
		Method: payment.Method{
			Type: "card",
			Card: payment.Card{
				Number:            "4111111111111111",
				Holder:            "John Doe",
				CVV:               "123",
				Expiration:        "12/2025",
				InstallmentNumber: 1,
			},
		},
	}
	createdPayment, _ := service.ProcessPayment(context.Background(), paymentRequest)

	result, err := service.RefundPayment(context.Background(), createdPayment.ID, 50.0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil || result.Status != "refunded" {
		t.Fatalf("expected payment to be refunded, got %v", result)
	}
}

func TestGetPayment_SuccessWithPrimaryProvider(t *testing.T) {
	primary := paypalProvider.New()
	secondary := stripeProvider.New()

	service := &payment.Service{
		PrimaryProvider:   primary,
		SecondaryProvider: secondary,
	}

	paymentRequest := &payment.Payment{
		Amount:   100.0,
		Currency: "USD",
		Method: payment.Method{
			Type: "card",
			Card: payment.Card{
				Number:            "4111111111111111",
				Holder:            "John Doe",
				CVV:               "123",
				Expiration:        "12/2025",
				InstallmentNumber: 1,
			},
		},
	}
	createdPayment, _ := service.ProcessPayment(context.Background(), paymentRequest)

	result, err := service.GetPayment(context.Background(), createdPayment.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil || result.ID != createdPayment.ID {
		t.Fatalf("expected payment ID to match, got %v", result)
	}
}

func TestRefundPayment_FailureWithBothProviders(t *testing.T) {
	primary := paypalProvider.New()
	secondary := stripeProvider.New()

	service := &payment.Service{
		PrimaryProvider:   primary,
		SecondaryProvider: secondary,
	}

	_, err := service.RefundPayment(context.Background(), "invalid-id", 50.0)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestGetPayment_FailureWithBothProviders(t *testing.T) {
	primary := paypalProvider.New()
	secondary := stripeProvider.New()

	service := &payment.Service{
		PrimaryProvider:   primary,
		SecondaryProvider: secondary,
	}

	_, err := service.GetPayment(context.Background(), "invalid-id")
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}
