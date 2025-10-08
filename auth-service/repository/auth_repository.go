package repository

import (
	"auth-service/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Create(user *model.User) (model.User, error)
	GetByEmail(string) (model.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) Create(user *model.User) (model.User, error) {
	result := r.db.Create(user)
	if result.Error != nil {
		return model.User{}, result.Error
	}
	return *user, nil
}

func (r *authRepository) GetByEmail(email string) (model.User, error) {
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil

}
