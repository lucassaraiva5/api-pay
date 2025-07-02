package adapters

import (
	"lucassaraiva5/api-pay/internal/app/adapters/handler"
	"lucassaraiva5/api-pay/internal/app/domain"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	payment *handler.PaymentHandler
}

func NewHandlers(services *domain.Services) *Handlers {
	return &Handlers{
		payment: handler.NewPaymentHandler(services),
	}
}

func (h *Handlers) Configure(server *echo.Echo) {
	h.payment.Configure(server)
}
