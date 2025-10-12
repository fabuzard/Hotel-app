package handler

import (
	"payment-service/dto"
	"payment-service/model"
	"payment-service/service"

	"github.com/labstack/echo/v4"
)

type PaymentHandler struct {
	service service.PaymentService
}

func NewPaymentHandler(s service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: s,
	}
}

// create payment
func (h *PaymentHandler) CreatePayment(c echo.Context) error {
	var req dto.CreatePaymentRequest
	// get user id from token
	authHeader := c.Request().Header.Get("Authorization")
	userIDFloat := c.Get("user_id").(float64)
	userID := int(userIDFloat)
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request"})
	}
	payment := model.Payment{
		BookingID: uint(req.BookingID),
	}
	paymentRes, err := h.service.CreatePayment(payment, userID, authHeader)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	return c.JSON(201, map[string]interface{}{"payment": paymentRes})
}

func (h *PaymentHandler) SimulateWebhook(c echo.Context) error {
	var req dto.XenditWebhookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request"})
	}
	authHeader := c.Request().Header.Get("Authorization")
	payment, err := h.service.SimulatePaymentWebhook(req.PaymentID, authHeader)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, map[string]interface{}{"message": "webhook procssed", "payment": payment})
}
