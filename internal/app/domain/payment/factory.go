package payment

func NewPayment(amount float64, currency string, description string, method Method) *Payment {
	return &Payment{
		Amount:      amount,
		Currency:    currency,
		Description: description,
		Method:      method,
	}
}

func NewMethod(paymentType string, card Card) Method {
	return Method{
		Type: paymentType,
		Card: card,
	}
}

func NewCard(number string, holder string, cvv string, expiration string, installmentNumber int) Card {
	return Card{
		Number:            number,
		Holder:            holder,
		CVV:               cvv,
		Expiration:        expiration,
		InstallmentNumber: installmentNumber,
	}
}

func NewService() *Service {
	return &Service{}
}
