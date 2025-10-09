package service

import (
	"reservation-service/dto"
	"reservation-service/model"
	"reservation-service/repository"
)

type RoomService interface {
	CreateRoom(dto.CreateRoomRequest) (model.Room, error)
	GetRoomByID(int) (model.Room, error)
	UpdateRoom(dto.UpdateRoomRequest, int) (model.Room, error)
	DeleteRoom(int) error
	ListRooms() ([]dto.RoomResponse, error)
}

type roomService struct {
	repo repository.RoomRepository
}

func NewRoomService(r repository.RoomRepository) RoomService {
	return &roomService{repo: r}
}

func (s *roomService) CreateRoom(req dto.CreateRoomRequest) (model.Room, error) {
	room := model.Room{
		RoomNumber:    req.RoomNumber,
		RoomType:      req.RoomType,
		PricePerNight: req.PricePerNight,
		MaxGuest:      req.MaxGuest,
		Status:        "available",
	}
	return s.repo.CreateRoom(&room)
}

func (s *roomService) GetRoomByID(id int) (model.Room, error) {
	return s.repo.GetRoomByID(string(id))
}

func (s *roomService) UpdateRoom(req dto.UpdateRoomRequest, id int) (model.Room, error) {
	room, err := s.repo.GetRoomByID(string(id))
	if err != nil {
		return model.Room{}, err
	}
	room.RoomNumber = req.RoomNumber
	room.RoomType = req.RoomType
	room.PricePerNight = req.PricePerNight
	room.MaxGuest = req.MaxGuest
	room.Status = req.Status
	return s.repo.UpdateRoom(&room)
}

func (s *roomService) DeleteRoom(id int) error {
	return s.repo.DeleteRoom(string(id))

}

func (s *roomService) ListRooms() ([]dto.RoomResponse, error) {
	rooms, err := s.repo.ListRooms()
	if err != nil {
		return nil, err
	}
	var roomResponses []dto.RoomResponse
	for _, room := range rooms {
		roomResponses = append(roomResponses, dto.RoomResponse{
			ID:            room.ID,
			RoomNumber:    room.RoomNumber,
			RoomType:      room.RoomType,
			PricePerNight: room.PricePerNight,
			MaxGuest:      room.MaxGuest,
			Status:        room.Status,
		})
	}
	return roomResponses, nil
}
