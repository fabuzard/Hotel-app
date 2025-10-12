package handler

import (
	"net/http"
	"reservation-service/dto"
	"reservation-service/model"
	"reservation-service/service"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type BookingHandler struct {
	service service.BookingService
}

func NewBookingHandler(s service.BookingService) *BookingHandler {
	return &BookingHandler{service: s}
}
func (h *BookingHandler) CreateBooking(c echo.Context) error {
	var req dto.CreateBooking
	// get user id from token
	userIDFloat := c.Get("user_id").(float64)
	userID := uint(userIDFloat)

	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request"})
	}

	checkinDate, _ := time.Parse("2006-01-02", req.CheckinDate)
	checkoutDate, _ := time.Parse("2006-01-02", req.CheckoutDate)
	// check if checkin is before checkout
	if !checkinDate.Before(checkoutDate) {
		return c.JSON(400, map[string]string{"error": "check-in date must be before check-out date"})
	}

	// check if checkin and checkout are not in the past
	today := time.Now().Truncate(24 * time.Hour)
	if checkinDate.Before(today) || checkoutDate.Before(today) {
		return c.JSON(400, map[string]string{"error": "check-in and check-out dates must not be in the past"})
	}

	booking := model.Booking{
		UserID:       userID,
		RoomID:       uint(req.RoomID),
		CheckinDate:  checkinDate,
		CheckoutDate: checkoutDate,
		Status:       "pending",
	}
	createdBooking, err := h.service.CreateBooking(booking)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	response := dto.BookingResponse{
		ID:           createdBooking.ID,
		UserID:       createdBooking.UserID,
		RoomID:       int(createdBooking.RoomID),
		CheckinDate:  createdBooking.CheckinDate.Format("2006-01-02"),
		CheckoutDate: createdBooking.CheckoutDate.Format("2006-01-02"),
		Status:       createdBooking.Status,
		TotalAmount:  createdBooking.TotalAmount,
	}
	return c.JSON(201, response)
}

// get booking by id
func (h *BookingHandler) GetBookingByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "invalid booking ID"})
	}
	booking, err := h.service.GetBookingByID(id)
	if err != nil {
		return c.JSON(404, map[string]string{"error": "booking not found"})
	}
	response := dto.BookingResponse{
		ID:           booking.ID,
		UserID:       booking.UserID,
		RoomID:       int(booking.RoomID),
		CheckinDate:  booking.CheckinDate.Format("2006-01-02"),
		CheckoutDate: booking.CheckoutDate.Format("2006-01-02"),
		Status:       booking.Status,
		TotalAmount:  booking.TotalAmount,
	}
	return c.JSON(200, response)

}
func (h *BookingHandler) UpdateBooking(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "invalid booking ID"})
	}
	var req dto.UpdateBooking
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request"})
	}
	updatedBooking, err := h.service.UpdateBooking(req, id)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	response := dto.BookingResponse{
		ID:           updatedBooking.ID,
		UserID:       updatedBooking.UserID,
		RoomID:       int(updatedBooking.RoomID),
		CheckinDate:  updatedBooking.CheckinDate.Format("2006-01-02"),
		CheckoutDate: updatedBooking.CheckoutDate.Format("2006-01-02"),
		Status:       updatedBooking.Status,
		TotalAmount:  updatedBooking.TotalAmount,
	}
	return c.JSON(200, response)
}

func (h *BookingHandler) WebhookUpdate(c echo.Context) error {
	// get user id from token
	userIDFloat := c.Get("user_id").(float64)
	userID := uint(userIDFloat)
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "invalid booking ID"})
	}
	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request"})
	}
	updatedBooking, err := h.service.UpdateWebhookStatus(id, req.Status, int(userID))
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	response := dto.BookingResponse{
		ID:           updatedBooking.ID,
		UserID:       updatedBooking.UserID,
		RoomID:       int(updatedBooking.RoomID),
		CheckinDate:  updatedBooking.CheckinDate.Format("2006-01-02"),
		CheckoutDate: updatedBooking.CheckoutDate.Format("2006-01-02"),
		Status:       updatedBooking.Status,
		TotalAmount:  updatedBooking.TotalAmount,
	}
	return c.JSON(200, response)
}

func (h *BookingHandler) DeleteBooking(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "invalid booking ID"})
	}
	err = h.service.DeleteBooking(id)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, map[string]string{"message": "booking deleted"})
}

func (h *BookingHandler) Checkin(c echo.Context) error {
	userIDFloat, ok := c.Get("user_id").(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"status":  "error",
			"message": "unauthorized: invalid token or missing user ID",
		})
	}
	userID := uint(userIDFloat)

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status":  "error",
			"message": "invalid booking ID",
		})
	}

	booking, err := h.service.Checkin(id, int(userID))
	if err != nil {
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "unauthorized") {
			statusCode = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "not available") {
			statusCode = http.StatusConflict
		}
		return c.JSON(statusCode, map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Check-in successful",
		"data":    booking,
	})
}

func (h *BookingHandler) Checkout(c echo.Context) error {
	userIDFloat, ok := c.Get("user_id").(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"status":  "error",
			"message": "unauthorized: invalid token or missing user ID",
		})
	}
	userID := uint(userIDFloat)

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status":  "error",
			"message": "invalid booking ID",
		})
	}

	booking, err := h.service.Checkout(id, int(userID))
	if err != nil {
		statusCode := http.StatusBadRequest
		if strings.Contains(err.Error(), "unauthorized") {
			statusCode = http.StatusUnauthorized
		}
		return c.JSON(statusCode, map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Checkout successful",
		"data":    booking,
	})
}
