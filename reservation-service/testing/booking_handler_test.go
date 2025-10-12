package testing

import (
	"net/http"
	"net/http/httptest"
	"reservation-service/dto"
	"reservation-service/handler"
	"reservation-service/model"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// ---- Dummy Booking Service ----
type dummyBookingService struct{}

func (d *dummyBookingService) CreateBooking(b model.Booking) (model.Booking, error) {
	b.ID = 1
	b.Status = "pending"
	b.TotalAmount = 500.0
	return b, nil
}

func (d *dummyBookingService) GetBookingByID(id int) (model.Booking, error) {
	return model.Booking{
		ID:           uint(id),
		UserID:       1,
		RoomID:       101,
		CheckinDate:  time.Now(),
		CheckoutDate: time.Now().Add(24 * time.Hour),
		Status:       "pending",
		TotalAmount:  500.0,
	}, nil
}

func (d *dummyBookingService) UpdateBooking(req dto.UpdateBooking, id int) (model.Booking, error) {
	return model.Booking{
		ID:           uint(id),
		UserID:       1,
		RoomID:       uint(req.RoomID),
		CheckinDate:  time.Now(),
		CheckoutDate: time.Now().Add(24 * time.Hour),
		Status:       req.Status,
		TotalAmount:  500.0,
	}, nil
}

func (d *dummyBookingService) DeleteBooking(id int) error {
	return nil
}

func (d *dummyBookingService) UpdateWebhookStatus(id int, status string, userID int) (model.Booking, error) {
	return model.Booking{
		ID:     uint(id),
		UserID: uint(userID),
		Status: status,
	}, nil
}

func (d *dummyBookingService) Checkin(id, userID int) (model.Booking, error) {
	return model.Booking{
		ID:     uint(id),
		UserID: uint(userID),
		Status: "checked_in",
	}, nil
}

func (d *dummyBookingService) Checkout(id, userID int) (model.Booking, error) {
	return model.Booking{
		ID:     uint(id),
		UserID: uint(userID),
		Status: "checked_out",
	}, nil
}
func (d *dummyBookingService) ListBookings() ([]dto.BookingResponse, error) {
	return []dto.BookingResponse{
		{
			ID:           1,
			UserID:       1,
			RoomID:       101,
			CheckinDate:  "2025-10-15",
			CheckoutDate: "2025-10-16",
			Status:       "pending",
			TotalAmount:  500.0,
		},
	}, nil
}

// ---- Tests ----
func TestCreateBooking_Success(t *testing.T) {
	e := echo.New()
	h := handler.NewBookingHandler(&dummyBookingService{})

	body := `{"room_id":101,"checkin_date":"2025-10-15","checkout_date":"2025-10-16"}`
	req := httptest.NewRequest(http.MethodPost, "/bookings", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", float64(1))

	_ = h.CreateBooking(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "pending")
	assert.Contains(t, rec.Body.String(), "101")
}

func TestGetBookingByID_Success(t *testing.T) {
	e := echo.New()
	h := handler.NewBookingHandler(&dummyBookingService{})

	req := httptest.NewRequest(http.MethodGet, "/bookings/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/bookings/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	_ = h.GetBookingByID(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "pending")
}

func TestCheckin_Success(t *testing.T) {
	e := echo.New()
	h := handler.NewBookingHandler(&dummyBookingService{})

	req := httptest.NewRequest(http.MethodPost, "/bookings/1/checkin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", float64(1))
	c.SetPath("/bookings/:id/checkin")
	c.SetParamNames("id")
	c.SetParamValues("1")

	_ = h.Checkin(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "checked_in")
}

func TestCheckout_Success(t *testing.T) {
	e := echo.New()
	h := handler.NewBookingHandler(&dummyBookingService{})

	req := httptest.NewRequest(http.MethodPost, "/bookings/1/checkout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", float64(1))
	c.SetPath("/bookings/:id/checkout")
	c.SetParamNames("id")
	c.SetParamValues("1")

	_ = h.Checkout(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "checked_out")
}
