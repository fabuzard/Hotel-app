package dto

// Create room request
type CreateRoomRequest struct {
	RoomNumber    int     `json:"room_number" validate:"required"`
	RoomType      string  `json:"room_type" validate:"required"`
	PricePerNight float64 `json:"price_per_night" validate:"required,gt=0"`
	MaxGuest      int     `json:"max_guest" validate:"required,gt=0"`
}

// Update room request
type UpdateRoomRequest struct {
	RoomNumber    int     `json:"room_number" validate:"required"`
	RoomType      string  `json:"room_type" validate:"required"`
	PricePerNight float64 `json:"price_per_night" validate:"required,gt=0"`
	MaxGuest      int     `json:"max_guest" validate:"required,gt=0"`
	Status        string  `json:"status" validate:"required,oneof=available booked maintenance"`
}

// Room response
type RoomResponse struct {
	ID            uint    `json:"id"`
	RoomNumber    int     `json:"room_number"`
	RoomType      string  `json:"room_type"`
	PricePerNight float64 `json:"price_per_night"`
	MaxGuest      int     `json:"max_guest"`
	Status        string  `json:"status"`
}

// List rooms response
type ListRoomsResponse struct {
	Rooms []RoomResponse `json:"rooms"`
}
