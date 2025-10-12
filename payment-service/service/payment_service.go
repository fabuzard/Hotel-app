package service

import (
	"fmt"
	"payment-service/client"
	"payment-service/model"
	"payment-service/repository"
)

type PaymentService interface {
	CreatePayment(payment model.Payment, userID int, authToken string) (string, error)
	GetPaymentByID(id int) (model.Payment, error)
	UpdatePayment(payment model.Payment) (model.Payment, error)
	DeletePayment(id int) error
	ListPaymentByUserID(userID int) ([]model.Payment, error)
	SimulatePaymentWebhook(paymentID int, authToken string) (model.Payment, error)
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(r repository.PaymentRepository) PaymentService {
	return &paymentService{
		repo: r,
	}
}

func (s *paymentService) CreatePayment(payment model.Payment, userID int, authToken string) (string, error) {
	// get the booking id from reservation service
	booking, err := client.GetBookingByID(payment.BookingID, authToken)
	if err != nil {
		return "", err
	}

	// validate the user id matches the booking user id
	if booking.UserID != uint(userID) {
		return "", err
	}

	// Validate booking status is pending and date is valid
	if booking.Status != "pending" {
		return "", err
	}
	// set amount from booking total amount
	payment.Amount = booking.TotalAmount

	// process payment with xendit
	url, err := client.CreateXenditPaymentURL(booking.ID, payment.Amount, "")
	if err != nil {
		return "", err
	}
	// set provider and provider payment id
	payment.Provider = "xendit"
	payment.ProviderPaymentID = "dummy-xendit-id"
	payment.Status = "pending"
	payment.PaymentURL = url
	// save payment record to db
	_, err = s.repo.CreatePayment(&payment)
	if err != nil {
		return "", err
	}

	// return payment URL from xendit
	return url, nil
}

// Simulate webhook
func (s *paymentService) SimulatePaymentWebhook(paymentID int, authToken string) (model.Payment, error) {
	payment, err := s.repo.GetPaymentByID(paymentID)
	if err != nil {
		return model.Payment{}, err
	}
	if payment.Status == "paid" {
		return model.Payment{}, fmt.Errorf("payment is already paid")
	}
	payment.Status = "paid"
	updatedPayment, err := s.repo.UpdatePayment(&payment)
	if err != nil {
		return model.Payment{}, err
	}

	// update booking status to confirmed
	err = client.UpdateBookingStatus(uint(updatedPayment.BookingID), "paid", authToken)
	if err != nil {
		return model.Payment{}, err
	}

	return updatedPayment, nil
}
func (s *paymentService) GetPaymentByID(id int) (model.Payment, error) {
	return s.repo.GetPaymentByID(id)
}
func (s *paymentService) UpdatePayment(payment model.Payment) (model.Payment, error) {
	return s.repo.UpdatePayment(&payment)
}
func (s *paymentService) DeletePayment(id int) error {
	return s.repo.DeletePayment(id)
}
func (s *paymentService) ListPaymentByUserID(userID int) ([]model.Payment, error) {
	return s.repo.ListPaymentByUserID(userID)
}
