package model

import "time"

// Rooms table
type Room struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	RoomNumber    int       `gorm:"size:50;not null" json:"room_number"`
	RoomType      string    `gorm:"size:100" json:"room_type"`
	PricePerNight float64   `gorm:"not null" json:"price_per_night"`
	MaxGuest      int       `json:"max_guest"`
	Status        string    `gorm:"type:varchar(20);default:'available'" json:"status"` // available, booked, maintenance
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Bookings table
type Booking struct {
	ID           string    `gorm:"primaryKey;autoIncrement"`
	UserID       string    `gorm:"type:uuid;not null" json:"user_id"`
	RoomID       string    `gorm:"type:uuid;not null" json:"room_id"`
	CheckinDate  time.Time `gorm:"type:date" json:"checkin_date"`
	CheckoutDate time.Time `gorm:"type:date" json:"checkout_date"`
	Status       string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, confirmed, checked_in, cancelled
	TotalAmount  float64   `json:"total_amount"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
