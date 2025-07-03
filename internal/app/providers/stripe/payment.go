package stripeProvider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lucassaraiva5/api-pay/internal/app/domain/stripe"
	"lucassaraiva5/api-pay/internal/app/transport/outbound"
	"lucassaraiva5/api-pay/internal/infra/logger"
	"lucassaraiva5/api-pay/internal/infra/logger/attributes"
	"net/http"
	"os"
)

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func getStripeMockURL() string {
	// Usa a variável de ambiente STRIPE_MOCK_URL (ex: http://stripe-mock:8082 no Docker Compose)
	url := os.Getenv("STRIPE_MOCK_URL")
	if url == "" {
		return "http://stripe-mock:8082" // padrão para desenvolvimento local
	}
	return url
}

func (p *Provider) CreatePayment(req interface{}) (interface{}, error) {
	ctx := context.Background()
	logger.Info(ctx, "[Stripe] CreatePayment called", nil)
	request, ok := req.(*stripe.PaymentRequest)
	if !ok {
		err := errors.New("invalid request type for Stripe")
		logger.Error(ctx, "[Stripe] Invalid request type", attributes.Attributes{"request": req}.WithError(err))
		return nil, err
	}

	payload := map[string]interface{}{
		"amount":              int64(request.Amount * 100),
		"currency":            request.Currency,
		"statementDescriptor": request.StatementDescriptor,
		"paymentType":         request.PaymentType,
		"description":         request.Description,
		"card":                request.Card,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		logger.Error(ctx, "[Stripe] Error marshaling payload", attributes.Attributes{"payload": payload}.WithError(err))
		return nil, err
	}
	url := getStripeMockURL() + "/transactions"
	logger.Info(ctx, "[Stripe] POST to mock", attributes.Attributes{"url": url})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Error(ctx, "[Stripe] Error on POST", attributes.Attributes{"url": url}.WithError(err))
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("Stripe mock returned status %d", resp.StatusCode)
		logger.Error(ctx, "[Stripe] Mock returned error status", attributes.Attributes{"status": resp.StatusCode, "url": url}.WithError(err))
		return nil, err
	}
	var mockResp struct {
		ID                  string      `json:"id"`
		Status              string      `json:"status"`
		Amount              int64       `json:"amount"`
		OriginalAmount      int64       `json:"originalAmount"`
		Currency            string      `json:"currency"`
		Description         string      `json:"description"`
		PaymentType         string      `json:"paymentType"`
		CreatedAt           string      `json:"date"`
		StatementDescriptor string      `json:"statementDescriptor"`
		Card                stripe.Card `json:"card"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&mockResp); err != nil {
		logger.Error(ctx, "[Stripe] Error decoding response", attributes.Attributes{}.WithError(err))
		return nil, err
	}
	payment := outbound.StripePaymentResponse{
		ID:          mockResp.ID,
		Status:      mockResp.Status,
		Amount:      float64(mockResp.Amount) / 100.0,
		Description: mockResp.Description,
		PaymentType: mockResp.PaymentType,
		CreatedAt:   mockResp.CreatedAt,
		Currency:    mockResp.Currency,
		Card:        mockResp.Card,
	}
	logger.Info(ctx, "[Stripe] Payment created successfully", attributes.Attributes{"payment_id": payment.ID})
	return &payment, nil
}

func (p *Provider) Refund(paymentID string) (interface{}, error) {
	ctx := context.Background()
	logger.Info(ctx, "[Stripe] Refund called", attributes.Attributes{"payment_id": paymentID})
	url := fmt.Sprintf("%s/void/%s", getStripeMockURL(), paymentID)
	logger.Info(ctx, "[Stripe] POST to mock for refund", attributes.Attributes{"url": url})
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		logger.Error(ctx, "[Stripe] Error on POST refund", attributes.Attributes{"url": url}.WithError(err))
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("Stripe mock returned status %d", resp.StatusCode)
		logger.Error(ctx, "[Stripe] Mock returned error status on refund", attributes.Attributes{"status": resp.StatusCode, "url": url}.WithError(err))
		return nil, err
	}
	var mockResp struct {
		ID                  string      `json:"id"`
		Status              string      `json:"status"`
		Amount              int64       `json:"amount"`
		OriginalAmount      int64       `json:"originalAmount"`
		Currency            string      `json:"currency"`
		Description         string      `json:"description"`
		PaymentType         string      `json:"paymentType"`
		CreatedAt           string      `json:"date"`
		StatementDescriptor string      `json:"statementDescriptor"`
		Card                stripe.Card `json:"card"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&mockResp); err != nil {
		logger.Error(ctx, "[Stripe] Error decoding refund response", attributes.Attributes{}.WithError(err))
		return nil, err
	}
	payment := outbound.StripePaymentResponse{
		ID:          mockResp.ID,
		Status:      mockResp.Status,
		Amount:      float64(mockResp.Amount) / 100.0,
		Description: mockResp.Description,
		PaymentType: mockResp.PaymentType,
		CreatedAt:   mockResp.CreatedAt,
		Currency:    mockResp.Currency,
		Card:        mockResp.Card,
	}
	logger.Info(ctx, "[Stripe] Refund processed successfully", attributes.Attributes{"payment_id": payment.ID})
	return &payment, nil
}

func (p *Provider) GetPayment(paymentID string) (interface{}, error) {
	fmt.Println("[DEBUG] STRIPE_MOCK_URL:", getStripeMockURL())
	url := fmt.Sprintf("%s/transactions/%s", getStripeMockURL(), paymentID)
	fmt.Println("[DEBUG] GET URL:", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Stripe mock returned status %d", resp.StatusCode)
	}
	var mockResp struct {
		ID                  string      `json:"id"`
		Status              string      `json:"status"`
		Amount              int64       `json:"amount"`
		OriginalAmount      int64       `json:"originalAmount"`
		Currency            string      `json:"currency"`
		Description         string      `json:"description"`
		PaymentType         string      `json:"paymentType"`
		CreatedAt           string      `json:"date"`
		StatementDescriptor string      `json:"statementDescriptor"`
		Card                stripe.Card `json:"card"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&mockResp); err != nil {
		return nil, err
	}
	payment := outbound.StripePaymentResponse{
		ID:          mockResp.ID,
		Status:      mockResp.Status,
		Amount:      float64(mockResp.Amount) / 100.0,
		Description: mockResp.Description,
		PaymentType: mockResp.PaymentType,
		CreatedAt:   mockResp.CreatedAt,
		Currency:    mockResp.Currency,
		Card:        mockResp.Card,
	}
	return &payment, nil
}
