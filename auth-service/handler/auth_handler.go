package handler

import (
	"auth-service/dto"
	"auth-service/service"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var user dto.RegisterRequest
	if err := c.Bind(&user); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request"})
	}
	createdUser, err := h.service.RegisterUser(user)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	response := dto.RegisterResponse{
		ID:       createdUser.ID,
		Email:    createdUser.Email,
		FullName: createdUser.FullName,
	}
	return c.JSON(201, response)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request"})
	}
	token, err := h.service.LoginRequest(req)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, map[string]string{"token": token})
}
