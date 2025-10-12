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
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	UserID       uint      `gorm:"not null"` // references auth_service.users.id
	RoomID       uint      `gorm:"not null"` // references rooms.id
	CheckinDate  time.Time `gorm:"type:date"`
	CheckoutDate time.Time `gorm:"type:date"`
	Status       string    `gorm:"type:varchar(20);default:'pending'"`
	TotalAmount  float64
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
