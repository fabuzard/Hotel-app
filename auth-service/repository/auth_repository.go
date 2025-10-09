package repository

import (
	"auth-service/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Create(user *model.Users) (model.Users, error)
	GetByEmail(string) (model.Users, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) Create(user *model.Users) (model.Users, error) {
	result := r.db.Create(user)
	if result.Error != nil {
		return model.Users{}, result.Error
	}

	return *user, nil
}

func (r *authRepository) GetByEmail(email string) (model.Users, error) {
	var user model.Users
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return model.Users{}, result.Error
	}
	return user, nil

}
