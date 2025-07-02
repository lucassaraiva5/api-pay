package payment

type Payment struct {
	ID                  string  `json:"id"`
	Amount              float64 `json:"amount"`
	Currency            string  `json:"currency"`
	Description         string  `json:"description"`
	Status              string  `json:"status"`
	CreatedAt           string  `json:"createdAt"`
	StatementDescriptor string  `json:"statementDescriptor,omitempty"`
	PaymentType         string  `json:"paymentType,omitempty"`
	CardID              string  `json:"cardId,omitempty"`
	Method              Method  `json:"method"`
}

type Method struct {
	Type string `json:"type"`
	Card Card   `json:"card"`
}

type Card struct {
	Number            string `json:"number"`
	Holder            string `json:"holder"`
	CVV               string `json:"cvv"`
	Expiration        string `json:"expiration"`
	InstallmentNumber int    `json:"installmentNumber"`
}

type Refund struct {
	Amount float64 `json:"amount"`
}
