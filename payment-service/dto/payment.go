package dto

// Payment request
type CreatePaymentRequest struct {
	BookingID int `json:"booking_id" validate:"required"`
}

// Payment response
type PaymentResponse struct {
	Url string `json:"url"`
}

// Xendit webhook request
type XenditWebhookRequest struct {
	PaymentID int `json:"payment_id" validate:"required"`
}
