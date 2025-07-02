package inbound

type (
	CreatePaymentRequest struct {
		ID          string  `json:"id"`
		Amount      float64 `json:"amount"`
		Currency    string  `json:"currency"`
		Description string  `json:"description"`
		Method      Method  `json:"method"`
		Status      string  `json:"status"` 
	}
)
