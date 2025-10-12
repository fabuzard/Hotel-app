package repository

import (
	"payment-service/model"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	// Define methods for payment repository
	CreatePayment(payment *model.Payment) (model.Payment, error)
	GetPaymentByID(id int) (model.Payment, error)
	UpdatePayment(payment *model.Payment) (model.Payment, error)
	DeletePayment(id int) error
	ListPaymentByUserID(userID int) ([]model.Payment, error)
}
type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) CreatePayment(payment *model.Payment) (model.Payment, error) {
	result := r.db.Create(payment)
	if result.Error != nil {
		return model.Payment{}, result.Error
	}
	return *payment, nil
}
func (r *paymentRepository) GetPaymentByID(id int) (model.Payment, error) {
	var payment model.Payment
	result := r.db.First(&payment, "id = ?", id)
	if result.Error != nil {
		return model.Payment{}, result.Error
	}
	return payment, nil
}
func (r *paymentRepository) UpdatePayment(payment *model.Payment) (model.Payment, error) {
	result := r.db.Save(payment)
	if result.Error != nil {
		return model.Payment{}, result.Error
	}
	return *payment, nil
}
func (r *paymentRepository) DeletePayment(id int) error {
	result := r.db.Delete(&model.Payment{}, "id = ?", id)
	return result.Error
}
func (r *paymentRepository) ListPaymentByUserID(userID int) ([]model.Payment, error) {
	var payments []model.Payment
	result := r.db.Where("user_id = ?", userID).Find(&payments)
	if result.Error != nil {
		return nil, result.Error
	}
	return payments, nil
}
