package paypalProvider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lucassaraiva5/api-pay/internal/app/domain/paypal"
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

func getPaypalMockURL() string {
	// Usa a variável de ambiente PAYPAL_MOCK_URL (ex: http://paypal-mock:8081 no Docker Compose)
	url := os.Getenv("PAYPAL_MOCK_URL")
	if url == "" {
		return "http://paypal-mock:8081" // padrão para desenvolvimento local
	}
	return url
}

func (p *Provider) CreatePayment(req interface{}) (interface{}, error) {
	ctx := context.Background()
	logger.Info(ctx, "[PayPal] CreatePayment called", nil)
	request, ok := req.(*paypal.PaymentRequest)
	if !ok {
		err := errors.New("invalid request type for PayPal")
		logger.Error(ctx, "[PayPal] Invalid request type", attributes.Attributes{"request": req}.WithError(err))
		return nil, err
	}

	payload := map[string]interface{}{
		"amount":        int64(request.Amount * 100), // assume amount is float in main currency, convert to cents
		"currency":      request.Currency,
		"description":   request.Description,
		"paymentMethod": request.PaymentMethod,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		logger.Error(ctx, "[PayPal] Error marshaling payload", attributes.Attributes{"payload": payload}.WithError(err))
		return nil, err
	}
	url := getPaypalMockURL() + "/charges"
	logger.Info(ctx, "[PayPal] POST to mock", attributes.Attributes{"url": url})
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Error(ctx, "[PayPal] Error on POST", attributes.Attributes{"url": url}.WithError(err))
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("PayPal mock returned status %d", resp.StatusCode)
		logger.Error(ctx, "[PayPal] Mock returned error status", attributes.Attributes{"status": resp.StatusCode, "url": url}.WithError(err))
		return nil, err
	}
	var mockResp struct {
		ID             string `json:"id"`
		Status         string `json:"status"`
		OriginalAmount int64  `json:"originalAmount"`
		CurrentAmount  int64  `json:"currentAmount"`
		Currency       string `json:"currency"`
		Description    string `json:"description"`
		CreatedAt      string `json:"createdAt"`
		PaymentMethod  string `json:"paymentMethod"`
		CardId         string `json:"cardId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&mockResp); err != nil {
		logger.Error(ctx, "[PayPal] Error decoding response", attributes.Attributes{}.WithError(err))
		return nil, err
	}
	payment := outbound.PaypalPaymentResponse{
		ID:          mockResp.ID,
		Status:      mockResp.Status,
		Amount:      float64(mockResp.CurrentAmount) / 100.0,
		Currency:    mockResp.Currency,
		Description: mockResp.Description,
		CreatedAt:   mockResp.CreatedAt,
		Method:      request.PaymentMethod,
	}
	logger.Info(ctx, "[PayPal] Payment created successfully", attributes.Attributes{"payment_id": payment.ID})
	return &payment, nil
}

func (p *Provider) Refund(paymentID string) (interface{}, error) {
	ctx := context.Background()
	logger.Info(ctx, "[PayPal] Refund called", attributes.Attributes{"payment_id": paymentID})
	url := fmt.Sprintf("%s/refund/%s", getPaypalMockURL(), paymentID)
	logger.Info(ctx, "[PayPal] POST to mock for refund", attributes.Attributes{"url": url})
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		logger.Error(ctx, "[PayPal] Error on POST refund", attributes.Attributes{"url": url}.WithError(err))
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("PayPal mock returned status %d", resp.StatusCode)
		logger.Error(ctx, "[PayPal] Mock returned error status on refund", attributes.Attributes{"status": resp.StatusCode, "url": url}.WithError(err))
		return nil, err
	}
	var mockResp struct {
		ID             string `json:"id"`
		Status         string `json:"status"`
		OriginalAmount int64  `json:"originalAmount"`
		CurrentAmount  int64  `json:"currentAmount"`
		Currency       string `json:"currency"`
		Description    string `json:"description"`
		CreatedAt      string `json:"createdAt"`
		PaymentMethod  string `json:"paymentMethod"`
		CardId         string `json:"cardId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&mockResp); err != nil {
		logger.Error(ctx, "[PayPal] Error decoding refund response", attributes.Attributes{}.WithError(err))
		return nil, err
	}
	payment := outbound.PaypalPaymentResponse{
		ID:          mockResp.ID,
		Status:      mockResp.Status,
		Amount:      float64(mockResp.CurrentAmount) / 100.0,
		Currency:    mockResp.Currency,
		Description: mockResp.Description,
		CreatedAt:   mockResp.CreatedAt,
		Method:      paypal.PaymentMethod{Type: mockResp.PaymentMethod},
	}
	logger.Info(ctx, "[PayPal] Refund processed successfully", attributes.Attributes{"payment_id": payment.ID})
	return &payment, nil
}

func (p *Provider) GetPayment(paymentID string) (interface{}, error) {
	fmt.Println("[DEBUG] PAYPAL_MOCK_URL:", getPaypalMockURL())
	url := fmt.Sprintf("%s/charges/%s", getPaypalMockURL(), paymentID)
	fmt.Println("[DEBUG] GET URL:", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PayPal mock returned status %d", resp.StatusCode)
	}
	var mockResp struct {
		ID             string `json:"id"`
		Status         string `json:"status"`
		OriginalAmount int64  `json:"originalAmount"`
		CurrentAmount  int64  `json:"currentAmount"`
		Currency       string `json:"currency"`
		Description    string `json:"description"`
		CreatedAt      string `json:"createdAt"`
		PaymentMethod  string `json:"paymentMethod"`
		CardId         string `json:"cardId"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&mockResp); err != nil {
		return nil, err
	}
	payment := outbound.PaypalPaymentResponse{
		ID:          mockResp.ID,
		Status:      mockResp.Status,
		Amount:      float64(mockResp.CurrentAmount) / 100.0,
		Currency:    mockResp.Currency,
		Description: mockResp.Description,
		CreatedAt:   mockResp.CreatedAt,
		Method:      paypal.PaymentMethod{Type: mockResp.PaymentMethod},
	}
	return &payment, nil
}
