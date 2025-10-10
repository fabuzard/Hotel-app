package repository

import (
	"reservation-service/model"

	"gorm.io/gorm"
)

type RoomRepository interface {
	CreateRoom(room *model.Room) (model.Room, error)
	GetRoomByID(id int) (model.Room, error)
	UpdateRoom(room *model.Room) (model.Room, error)
	DeleteRoom(id string) error
	ListRooms() ([]model.Room, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db}
}
func (r *roomRepository) CreateRoom(room *model.Room) (model.Room, error) {
	result := r.db.Create(room)
	if result.Error != nil {
		return model.Room{}, result.Error
	}
	return *room, nil
}

func (r *roomRepository) GetRoomByID(id int) (model.Room, error) {
	var room model.Room
	result := r.db.First(&room, "id = ?", id)
	if result.Error != nil {
		return model.Room{}, result.Error
	}
	return room, nil
}
func (r *roomRepository) UpdateRoom(room *model.Room) (model.Room, error) {
	result := r.db.Save(room)
	if result.Error != nil {
		return model.Room{}, result.Error
	}
	return *room, nil
}
func (r *roomRepository) DeleteRoom(id string) error {
	result := r.db.Delete(&model.Room{}, "id = ?", id)
	return result.Error
}

func (r *roomRepository) ListRooms() ([]model.Room, error) {
	var rooms []model.Room
	result := r.db.Find(&rooms)
	if result.Error != nil {
		return nil, result.Error
	}
	return rooms, nil
}
