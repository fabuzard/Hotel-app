package worker

import (
	"log"
	"time"

	"reservation-service/model"

	"gorm.io/gorm"
)

func StartScheduler(db *gorm.DB) {
	log.Println("🏃 Running daily checkout worker...")

	today := time.Now().Truncate(24 * time.Hour)
	dateStr := today.Format("2006-01-02")

	var bookings []model.Booking

	// Find bookings with checkout date = today and status = "checked_in"
	err := db.Where("DATE(checkout_date) = ? AND status = ?", dateStr, "checked_in").Find(&bookings).Error
	if err != nil {
		log.Println("❌ Failed to get bookings:", err)
		return
	}

	if len(bookings) == 0 {
		log.Println("✅ No bookings to checkout today")
		return
	}

	for _, b := range bookings {
		// Mark booking as completed
		if err := db.Model(&model.Booking{}).Where("id = ?", b.ID).
			Update("status", "completed").Error; err != nil {
			log.Printf("❌ Failed to update booking %d: %v\n", b.ID, err)
			continue
		}

		// Mark room as available
		if err := db.Model(&model.Room{}).Where("id = ?", b.RoomID).
			Update("available", true).Error; err != nil {
			log.Printf("❌ Failed to mark room %d available: %v\n", b.RoomID, err)
			continue
		}

		log.Printf("✅ Checked out booking %d and freed room %d\n", b.ID, b.RoomID)
	}
}
