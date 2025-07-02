package inbound

type Card struct {
	Number            string `json:"number"`
	Holder            string `json:"holder"`
	CVV               string `json:"cvv"`
	Expiration        string `json:"expiration"`
	InstallmentNumber int    `json:"installmentNumber"`
}
