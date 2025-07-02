package paypal

type PaymentMethod struct {
	Type string `json:"type"`
	Card Card   `json:"card"`
}

type Card struct {
	Number       string `json:"number"`
	HolderName   string `json:"holderName"`
	CVV          string `json:"cvv"`
	Expiration   string `json:"expirationDate"`
	Installments int    `json:"installments"`
}

type PaymentRequest struct {
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	Description   string        `json:"description"`
	PaymentMethod PaymentMethod `json:"paymentMethod"`
}
