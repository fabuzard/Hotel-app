package service

import (
	"fmt"
	"reservation-service/dto"
	"reservation-service/model"
	"reservation-service/repository"
)

type BookingService interface {
	CreateBooking(booking model.Booking) (model.Booking, error)
	GetBookingByID(int) (model.Booking, error)
	UpdateBooking(dto.UpdateBooking, int) (model.Booking, error)
	DeleteBooking(int) error
	ListBookings() ([]dto.BookingResponse, error)
	UpdateWebhookStatus(id int, status string, userID int) (model.Booking, error)
}

type bookingService struct {
	repo repository.BookingRepository
}

func NewBookingService(r repository.BookingRepository) BookingService {
	return &bookingService{
		repo: r,
	}
}

func (s *bookingService) CreateBooking(booking model.Booking) (model.Booking, error) {
	// Get room info
	room, err := s.repo.GetRoomByID(int(booking.RoomID))
	if err != nil {
		return model.Booking{}, err
	}
	// check if room is available on the requested dates
	// get all bookings for the room
	bookings, err := s.repo.GetBookingsByRoomID(int(booking.RoomID))
	if err != nil {
		return model.Booking{}, err
	}
	for _, b := range bookings {
		// check if the requested dates overlap with existing bookings
		if (booking.CheckinDate.Before(b.CheckoutDate) && booking.CheckinDate.After(b.CheckinDate)) ||
			(booking.CheckoutDate.After(b.CheckinDate) && booking.CheckoutDate.Before(b.CheckoutDate)) ||
			(booking.CheckinDate.Equal(b.CheckinDate) || booking.CheckoutDate.Equal(b.CheckoutDate)) {
			return model.Booking{}, fmt.Errorf("room is not available on the requested dates")
		}
	}
	// if not available, return error
	// calculate number of nights
	nights := int(booking.CheckoutDate.Sub(booking.CheckinDate).Hours() / 24)
	if nights <= 0 {
		return model.Booking{}, fmt.Errorf("invalid check-in and check-out dates")
	}
	// calculate total amount
	booking.TotalAmount = float64(nights) * room.PricePerNight
	// save booking
	return s.repo.CreateBooking(&booking)

}

func (s *bookingService) GetBookingByID(id int) (model.Booking, error) {
	return s.repo.GetBookingByID(id)
}
func (s *bookingService) UpdateBooking(req dto.UpdateBooking, id int) (model.Booking, error) {
	booking, err := s.repo.GetBookingByID(id)
	if err != nil {
		return model.Booking{}, err
	}

	return s.repo.UpdateBooking(&booking)
}
func (s *bookingService) DeleteBooking(id int) error {
	return s.repo.DeleteBooking(id)
}
func (s *bookingService) ListBookings() ([]dto.BookingResponse, error) {
	bookings, err := s.repo.ListBookings()
	if err != nil {
		return nil, err
	}
	var bookingResponses []dto.BookingResponse
	for _, b := range bookings {
		bookingResponses = append(bookingResponses, dto.BookingResponse{
			ID:           b.ID,
			UserID:       b.UserID,
			RoomID:       int(b.RoomID),
			CheckinDate:  b.CheckinDate.Format("2006-01-02"),
			CheckoutDate: b.CheckoutDate.Format("2006-01-02"),
			Status:       b.Status,
			TotalAmount:  b.TotalAmount,
		})
	}
	return bookingResponses, nil
}

func (s *bookingService) UpdateWebhookStatus(id int, status string, userID int) (model.Booking, error) {
	// update booking status
	booking, err := s.repo.GetBookingByID(id)
	if err != nil {
		return model.Booking{}, err
	}
	if booking.UserID != uint(userID) {
		return model.Booking{}, fmt.Errorf("unauthorized you can only update your own booking")
	}
	// check if booking status is confirmed
	if booking.Status == "confirmed" {
		return model.Booking{}, fmt.Errorf("booking is already confirmed, you can now check-in")
	}
	booking.Status = status

	return s.repo.UpdateBooking(&booking)
}
