package stripeProvider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"lucassaraiva5/api-pay/internal/app/domain/stripe"
	"lucassaraiva5/api-pay/internal/app/transport/outbound"
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
	fmt.Println("[DEBUG] STRIPE_MOCK_URL:", getStripeMockURL())
	request, ok := req.(*stripe.PaymentRequest)
	if !ok {
		return nil, errors.New("invalid request type for Stripe")
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
		return nil, err
	}
	url := getStripeMockURL() + "/transactions"
	fmt.Println("[DEBUG] POST URL:", url)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
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

func (p *Provider) Refund(paymentID string) (interface{}, error) {
	fmt.Println("[DEBUG] STRIPE_MOCK_URL:", getStripeMockURL())
	url := fmt.Sprintf("%s/void/%s", getStripeMockURL(), paymentID)
	fmt.Println("[DEBUG] POST URL:", url)
	resp, err := http.Post(url, "application/json", nil)
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
