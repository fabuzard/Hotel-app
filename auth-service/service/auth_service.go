package service

import (
	"auth-service/dto"
	"auth-service/model"
	"auth-service/repository"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(dto.RegisterRequest) (model.User, error)
	LoginRequest(dto.LoginRequest) (model.User, error)
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) AuthService {
	return &authService{repo: r}
}

func (s *authService) RegisterUser(req dto.RegisterRequest) (model.User, error) {
	// check if already registered
	_, err := s.repo.GetByEmail(req.Email)
	if err == nil {
		return model.User{}, fmt.Errorf("email already registered")
	}

	if req.Password == "" {
		return model.User{}, fmt.Errorf("password cannot be empty")
	}
	// hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	req.Password = string(hashedBytes)
	// create user
	user := model.User{
		Email:        req.Email,
		PasswordHash: req.Password,
		FullName:     req.FullName,
	}
	return s.repo.Create(&user)
}

func (s *authService) LoginRequest(req dto.LoginRequest) (model.User, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return model.User{}, fmt.Errorf("invalid email or password")
	}
	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {

		return model.User{}, fmt.Errorf("invalid email or password")
	}
	return user, nil
}
