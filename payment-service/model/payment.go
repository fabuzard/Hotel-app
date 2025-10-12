package model

import (
	"time"
)

type Payment struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID         uint      `gorm:"not null" json:"booking_id"`       // FK to bookings.id
	Provider          string    `gorm:"size:50;not null" json:"provider"` // e.g. "xendit", "midtrans"
	ProviderPaymentID string    `gorm:"size:100" json:"provider_payment_id"`
	Amount            float64   `gorm:"not null" json:"amount"`
	Status            string    `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, paid, failed, refunded
	PaymentURL        string    `gorm:"size:255" json:"payment_url"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
