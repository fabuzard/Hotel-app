package dto

// Create booking request
type CreateBooking struct {
	RoomID       int    `json:"room_id" validate:"required"`
	CheckinDate  string `json:"checkin_date" validate:"required,datetime=2006-01-02"`
	CheckoutDate string `json:"checkout_date" validate:"required,datetime=2006-01-02,gtfield=CheckinDate"`
}

// Update booking request
type UpdateBooking struct {
	UserID       uint   `json:"user_id" validate:"required"`
	RoomID       int    `json:"room_id" validate:"required"`
	CheckinDate  string `json:"checkin_date" validate:"required,datetime=2006-01-02"`
	CheckoutDate string `json:"checkout_date" validate:"required,datetime=2006-01-02,gtfield=CheckinDate"`
	Status       string `json:"status" validate:"required,oneof=pending confirmed checked_in cancelled"`
}

// Booking response
type BookingResponse struct {
	ID           uint    `json:"id"`
	UserID       uint    `json:"user_id"`
	RoomID       int     `json:"room_id"`
	CheckinDate  string  `json:"checkin_date"`
	CheckoutDate string  `json:"checkout_date"`
	Status       string  `json:"status"`
	TotalAmount  float64 `json:"total_amount"`
}

// List bookings response
type ListBookingsResponse struct {
	Bookings []BookingResponse `json:"bookings"`
}
