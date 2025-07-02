package stripe

type Card struct {
	Number            string `json:"number"`
	Holder            string `json:"holder"`
	CVV               string `json:"cvv"`
	Expiration        string `json:"expiration"`
	InstallmentNumber int    `json:"installmentNumber"`
}

type PaymentRequest struct {
	Amount              float64 `json:"amount"`
	Currency            string  `json:"currency"`
	StatementDescriptor string  `json:"statementDescriptor"`
	PaymentType         string  `json:"paymentType"`
	Description         string  `json:"description"`
	Card                Card    `json:"card"`
}
