package handler

import (
	"reservation-service/dto"
	"reservation-service/service"

	"github.com/labstack/echo/v4"
)

type RoomHandler struct {
	service service.RoomService
}

func NewRoomHandler(s service.RoomService) *RoomHandler {
	return &RoomHandler{service: s}
}

func (h *RoomHandler) CreateRoom(c echo.Context) error {
	var req dto.CreateRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request"})
	}
	createdRoom, err := h.service.CreateRoom(req)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	response := dto.RoomResponse{
		ID:            createdRoom.ID,
		RoomNumber:    createdRoom.RoomNumber,
		RoomType:      createdRoom.RoomType,
		PricePerNight: createdRoom.PricePerNight,
		MaxGuest:      createdRoom.MaxGuest,
		Status:        createdRoom.Status,
	}
	return c.JSON(201, response)
}
