package repository

import (
	"reservation-service/model"

	"gorm.io/gorm"
)

type BookingRepository interface {
	// Define methods for booking repository
	CreateBooking(booking *model.Booking) (model.Booking, error)
	GetBookingByID(id int) (model.Booking, error)
	UpdateBookingByID(id int) (model.Booking, error)
	UpdateBooking(booking *model.Booking) (model.Booking, error)
	DeleteBooking(id int) error
	ListBookings() ([]model.Booking, error)
	GetBookingsByRoomID(roomID int) ([]model.Booking, error)
	GetRoomByID(id int) (model.Room, error)
	UpdateRoomStatus(roomID int, status string) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db}
}

func (r *bookingRepository) CreateBooking(booking *model.Booking) (model.Booking, error) {
	result := r.db.Create(booking)
	if result.Error != nil {
		return model.Booking{}, result.Error
	}
	return *booking, nil
}
func (r *bookingRepository) GetBookingByID(id int) (model.Booking, error) {
	var booking model.Booking
	result := r.db.First(&booking, "id = ?", id)
	if result.Error != nil {
		return model.Booking{}, result.Error
	}
	return booking, nil
}

func (r *bookingRepository) UpdateBooking(booking *model.Booking) (model.Booking, error) {
	result := r.db.Save(booking)
	if result.Error != nil {
		return model.Booking{}, result.Error
	}
	return *booking, nil
}
func (r *bookingRepository) DeleteBooking(id int) error {
	result := r.db.Delete(&model.Booking{}, "id = ?", id)
	return result.Error
}
func (r *bookingRepository) ListBookings() ([]model.Booking, error) {
	var bookings []model.Booking
	result := r.db.Find(&bookings)
	if result.Error != nil {
		return nil, result.Error
	}
	return bookings, nil
}
func (r *bookingRepository) GetBookingsByRoomID(roomID int) ([]model.Booking, error) {
	var bookings []model.Booking
	result := r.db.Where("room_id = ?", roomID).Find(&bookings)
	if result.Error != nil {
		return nil, result.Error
	}
	return bookings, nil
}
func (r *bookingRepository) GetRoomByID(id int) (model.Room, error) {
	var room model.Room
	result := r.db.First(&room, "id = ?", id)
	if result.Error != nil {
		return model.Room{}, result.Error
	}
	return room, nil
}

func (r *bookingRepository) UpdateBookingByID(id int) (model.Booking, error) {
	var booking model.Booking
	result := r.db.First(&booking, "id = ?", id)
	if result.Error != nil {
		return model.Booking{}, result.Error
	}
	// update status
	booking.Status = "confirmed"
	result = r.db.Save(&booking)
	if result.Error != nil {
		return model.Booking{}, result.Error
	}
	return booking, nil
}
func (r *bookingRepository) UpdateRoomStatus(roomID int, status string) error {
	var room model.Room
	result := r.db.First(&room, "id = ?", roomID)
	if result.Error != nil {
		return result.Error
	}
	room.Status = status
	result = r.db.Save(&room)
	return result.Error
}
