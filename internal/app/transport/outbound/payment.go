package outbound

type (
	CreatePaymentResponse struct {
		Id     string `json:"id"`
		Status string `json:"status"`
	}

	ReadPaymentByIdResponse struct {
		Id     string `json:"id"`
		Status string `json:"status"` 
	}

	RefundResponse struct {
		Id     string `json:"id"`
		Status string `json:"status"` 
	}
)
