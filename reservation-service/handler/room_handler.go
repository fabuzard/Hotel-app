package handler

import (
	"reservation-service/dto"
	"reservation-service/service"
	"strconv"

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

// get list of rooms
func (h *RoomHandler) ListRooms(c echo.Context) error {
	rooms, err := h.service.ListRooms()
	if err != nil {
		return c.JSON(500, map[string]string{"error": "failed to fetch rooms"})
	}
	return c.JSON(200, rooms)
}

func (h *RoomHandler) GetRoomByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "invalid room ID"})
	}
	room, err := h.service.GetRoomByID(id)
	if err != nil {
		return c.JSON(404, map[string]string{"error": "room not found"})
	}
	response := dto.RoomResponse{
		ID:            room.ID,
		RoomNumber:    room.RoomNumber,
		RoomType:      room.RoomType,
		PricePerNight: room.PricePerNight,
		MaxGuest:      room.MaxGuest,
		Status:        room.Status,
	}
	return c.JSON(200, response)
}

func (h *RoomHandler) UpdateRoom(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "invalid room ID"})
	}
	var req dto.UpdateRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request"})
	}
	updatedRoom, err := h.service.UpdateRoom(req, id)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	response := dto.RoomResponse{
		ID:            updatedRoom.ID,
		RoomNumber:    updatedRoom.RoomNumber,
		RoomType:      updatedRoom.RoomType,
		PricePerNight: updatedRoom.PricePerNight,
		MaxGuest:      updatedRoom.MaxGuest,
		Status:        updatedRoom.Status,
	}
	return c.JSON(200, response)
}
