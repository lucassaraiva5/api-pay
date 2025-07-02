package handler

import (
	"lucassaraiva5/api-pay/internal/app/domain"
	"net/http"

	"github.com/labstack/echo/v4"
	"lucassaraiva5/api-pay/internal/app/domain/payment"
	"lucassaraiva5/api-pay/internal/app/transport/inbound"
	"lucassaraiva5/api-pay/internal/app/transport/mapper"
)

type PaymentHandler struct {
	service *payment.Service
}

func NewPaymentHandler(services *domain.Services) *PaymentHandler {
	return &PaymentHandler{service: services.PaymentService}
}

func (h *PaymentHandler) Configure(server *echo.Echo) {
	server.POST("/payments", h.CreatePayment)
	server.POST("/refunds", h.RefundPayment)
	server.GET("/payments/:id", h.GetPayment)
}

func (h *PaymentHandler) CreatePayment(c echo.Context) error {
	var request inbound.CreatePaymentRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	payment := mapper.PaymentFromCreatePaymentRequest(&request)
	result, err := h.service.ProcessPayment(c.Request().Context(), payment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *PaymentHandler) RefundPayment(c echo.Context) error {
	var refund payment.Refund
	if err := c.Bind(&refund); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	paymentID := c.QueryParam("id")
	result, err := h.service.RefundPayment(c.Request().Context(), paymentID, refund.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func (h *PaymentHandler) GetPayment(c echo.Context) error {
	paymentID := c.Param("id")
	result, err := h.service.GetPayment(c.Request().Context(), paymentID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
