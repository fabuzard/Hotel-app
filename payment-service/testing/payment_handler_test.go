package testing

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"payment-service/handler"
	"payment-service/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// ---- Dummy service ----
type dummyService struct {
	shouldFail bool
}

func (d *dummyService) CreatePayment(payment model.Payment, userID int, auth string) (string, error) {
	if d.shouldFail {
		return "", errors.New("failed to create payment")
	}
	return "https://dummy-payment-url", nil
}

func (d *dummyService) SimulatePaymentWebhook(paymentID int, auth string) (model.Payment, error) {
	if d.shouldFail {
		return model.Payment{}, errors.New("failed webhook")
	}
	return model.Payment{
		ID:        uint(paymentID),
		Status:    "paid",
		BookingID: 10,
	}, nil
}

func (d *dummyService) DeletePayment(id int) error                   { return nil }
func (d *dummyService) GetPaymentByID(id int) (model.Payment, error) { return model.Payment{}, nil }
func (d *dummyService) UpdatePayment(p model.Payment) (model.Payment, error) {
	return p, nil
}
func (d *dummyService) ListPaymentByUserID(uid int) ([]model.Payment, error) {
	return nil, nil
}

// ---- Test: Success case ----
func TestCreatePayment_Success(t *testing.T) {
	e := echo.New()
	h := handler.NewPaymentHandler(&dummyService{})

	req := httptest.NewRequest(http.MethodPost, "/payments", strings.NewReader(`{"booking_id":1}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", float64(1))

	err := h.CreatePayment(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "payment_url")
}

// ---- Test: Failure case ----
func TestCreatePayment_Failure(t *testing.T) {
	e := echo.New()
	h := handler.NewPaymentHandler(&dummyService{shouldFail: true})

	req := httptest.NewRequest(http.MethodPost, "/payments", strings.NewReader(`{"booking_id":1}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", float64(1))

	err := h.CreatePayment(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to create payment")
}

// Test: SimulateWebhook ----
func TestSimulateWebhook_Success(t *testing.T) {
	e := echo.New()
	h := handler.NewPaymentHandler(&dummyService{})

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(`{"payment_id": 123}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer dummy_token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.SimulateWebhook(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"webhook procssed"`)
	assert.Contains(t, rec.Body.String(), `"payment"`)
}

func TestSimulateWebhook_Failure(t *testing.T) {
	e := echo.New()
	h := handler.NewPaymentHandler(&dummyService{shouldFail: true})

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(`{"payment_id": 123}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer dummy_token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.SimulateWebhook(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed webhook")
}
