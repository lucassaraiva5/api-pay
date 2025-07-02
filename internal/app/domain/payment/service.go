package payment

import (
	"context"
	"errors"
	"lucassaraiva5/api-pay/internal/app/domain/paypal"
	"lucassaraiva5/api-pay/internal/app/domain/stripe"
	"lucassaraiva5/api-pay/internal/app/providers/paypal"
	"lucassaraiva5/api-pay/internal/app/providers/stripe"
	"lucassaraiva5/api-pay/internal/app/transport/outbound"
)

type Service struct {
	PrimaryProvider   Provider
	SecondaryProvider Provider
}

func New(paypalProvider *paypalProvider.Provider, stripeProvider *stripeProvider.Provider) *Service {
	return &Service{
		PrimaryProvider:   paypalProvider,
		SecondaryProvider: stripeProvider,
	}
}

type Provider interface {
	CreatePayment(req interface{}) (interface{}, error)
	Refund(paymentID string, amount float64) (interface{}, error)
	GetPayment(paymentID string) (interface{}, error)
}

func (s *Service) ProcessPayment(ctx context.Context, payment *Payment) (*Payment, error) {
	result, err := s.processWithProvider(s.PrimaryProvider, payment)
	if err == nil {
		return result, nil
	}

	result, err = s.processWithProvider(s.SecondaryProvider, payment)
	if err != nil {
		return nil, errors.New("both providers failed: " + err.Error())
	}

	return result, nil
}

func (s *Service) RefundPayment(ctx context.Context, paymentID string, amount float64) (*Payment, error) {
	result, err := s.refundWithProvider(s.PrimaryProvider, paymentID, amount)
	if err == nil {
		return result, nil
	}

	result, err = s.refundWithProvider(s.SecondaryProvider, paymentID, amount)
	if err != nil {
		return nil, errors.New("both providers failed to refund: " + err.Error())
	}

	return result, nil
}

func (s *Service) GetPayment(ctx context.Context, paymentID string) (*Payment, error) {
	result, err := s.getWithProvider(s.PrimaryProvider, paymentID)
	if err == nil {
		return result, nil
	}

	result, err = s.getWithProvider(s.SecondaryProvider, paymentID)
	if err != nil {
		return nil, errors.New("both providers failed to retrieve payment: " + err.Error())
	}

	return result, nil
}

func (s *Service) processWithProvider(provider Provider, payment *Payment) (*Payment, error) {
	if provider == nil {
		return nil, errors.New("provider is nil")
	}

	var request interface{}

	switch provider.(type) {
	case *paypalProvider.Provider:
		request = &paypal.PaymentRequest{
			Amount:      payment.Amount,
			Currency:    payment.Currency,
			Description: payment.Description,
			PaymentMethod: paypal.PaymentMethod{
				Type: "card",
				Card: paypal.Card{
					Number:       payment.Method.Card.Number,
					HolderName:   payment.Method.Card.Holder,
					CVV:          payment.Method.Card.CVV,
					Expiration:   payment.Method.Card.Expiration,
					Installments: payment.Method.Card.InstallmentNumber,
				},
			},
		}
	case *stripeProvider.Provider:
		request = &stripe.PaymentRequest{
			Amount:              payment.Amount,
			Currency:            payment.Currency,
			StatementDescriptor: payment.Description,
			PaymentType:         "card",
			Card: stripe.Card{
				Number:            payment.Method.Card.Number,
				Holder:            payment.Method.Card.Holder,
				CVV:               payment.Method.Card.CVV,
				Expiration:        payment.Method.Card.Expiration,
				InstallmentNumber: payment.Method.Card.InstallmentNumber,
			},
		}
	default:
		return nil, errors.New("unsupported provider type")
	}

	response, err := provider.CreatePayment(request)
	if err != nil {
		return nil, err
	}

	return s.adaptResponse(response)
}

func (s *Service) refundWithProvider(provider Provider, paymentID string, amount float64) (*Payment, error) {
	if provider == nil {
		return nil, errors.New("provider is nil")
	}

	response, err := provider.Refund(paymentID, amount)
	if err != nil {
		return nil, err
	}

	return s.adaptResponse(response)
}

func (s *Service) getWithProvider(provider Provider, paymentID string) (*Payment, error) {
	if provider == nil {
		return nil, errors.New("provider is nil")
	}

	response, err := provider.GetPayment(paymentID)
	if err != nil {
		return nil, err
	}

	return s.adaptResponse(response)
}

func (s *Service) adaptResponse(response interface{}) (*Payment, error) {
	switch res := response.(type) {
	case *outbound.PaypalPaymentResponse:
		return &Payment{
			ID:          res.ID,
			Amount:      res.Amount,
			Currency:    res.Currency,
			Status:      res.Status,
			Description: res.Description,
			CreatedAt:   res.CreatedAt,
			Method: Method{
				Type: res.Method.Type,
				Card: Card{
					Number:            res.Method.Card.Number,
					Holder:            res.Method.Card.HolderName,
					CVV:               res.Method.Card.CVV,
					Expiration:        res.Method.Card.Expiration,
					InstallmentNumber: res.Method.Card.Installments,
				},
			},
		}, nil
	case *outbound.StripePaymentResponse:
		return &Payment{
			ID:          res.ID,
			Amount:      res.Amount,
			Currency:    res.Currency,
			Status:      res.Status,
			Description: res.Description,
			CreatedAt:   res.CreatedAt,
			Method: Method{
				Type: res.PaymentType,
				Card: Card{
					Number:            res.Card.Number,
					Holder:            res.Card.Holder,
					CVV:               res.Card.CVV,
					Expiration:        res.Card.Expiration,
					InstallmentNumber: res.Card.InstallmentNumber,
				},
			},
		}, nil
	default:
		return nil, errors.New("unsupported response type")
	}
}
