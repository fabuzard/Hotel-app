package dto

type BookingResponse struct {
	ID           uint    `json:"id"`
	UserID       uint    `json:"user_id"`
	RoomID       int     `json:"room_id"`
	CheckinDate  string  `json:"checkin_date"`
	CheckoutDate string  `json:"checkout_date"`
	Status       string  `json:"status"`
	TotalAmount  float64 `json:"total_amount"`
}
