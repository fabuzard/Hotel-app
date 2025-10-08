package model

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	Email        string    `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	FullName     string    `gorm:"size:255"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
