package paypalProvider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"lucassaraiva5/api-pay/internal/app/domain/paypal"
	"lucassaraiva5/api-pay/internal/app/transport/outbound"
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
	fmt.Println("[DEBUG] PAYPAL_MOCK_URL:", getPaypalMockURL())
	request, ok := req.(*paypal.PaymentRequest)
	if !ok {
		return nil, errors.New("invalid request type for PayPal")
	}

	payload := map[string]interface{}{
		"amount":        int64(request.Amount * 100), // assume amount is float in main currency, convert to cents
		"currency":      request.Currency,
		"description":   request.Description,
		"paymentMethod": request.PaymentMethod,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	url := getPaypalMockURL() + "/charges"
	fmt.Println("[DEBUG] POST URL:", url)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
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
		Method:      request.PaymentMethod,
	}
	return &payment, nil
}

func (p *Provider) Refund(paymentID string, amount float64) (interface{}, error) {
	fmt.Println("[DEBUG] PAYPAL_MOCK_URL:", getPaypalMockURL())
	payload := map[string]interface{}{
		"amount": int64(amount * 100),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/refund/%s", getPaypalMockURL(), paymentID)
	fmt.Println("[DEBUG] POST URL:", url)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
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
