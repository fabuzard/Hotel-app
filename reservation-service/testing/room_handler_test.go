package testing

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"reservation-service/dto"
	"reservation-service/handler"
	"reservation-service/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// ---- Dummy Room Service ----
type dummyRoomService struct{}

func (d *dummyRoomService) CreateRoom(req dto.CreateRoomRequest) (model.Room, error) {
	return model.Room{
		ID:            1,
		RoomNumber:    req.RoomNumber,
		RoomType:      req.RoomType,
		PricePerNight: req.PricePerNight,
		MaxGuest:      req.MaxGuest,
		Status:        "available",
	}, nil
}

func (d *dummyRoomService) ListRooms() ([]dto.RoomResponse, error) {
	return []dto.RoomResponse{
		{
			ID:            1,
			RoomNumber:    101,
			RoomType:      "Deluxe",
			PricePerNight: 100.0,
			MaxGuest:      2,
			Status:        "available",
		},
	}, nil
}

func (d *dummyRoomService) GetRoomByID(id int) (model.Room, error) {
	return model.Room{
		ID:            uint(id),
		RoomNumber:    101,
		RoomType:      "Deluxe",
		PricePerNight: 100.0,
		MaxGuest:      2,
		Status:        "available",
	}, nil
}

func (d *dummyRoomService) UpdateRoom(req dto.UpdateRoomRequest, id int) (model.Room, error) {
	return model.Room{
		ID:            uint(id),
		RoomNumber:    req.RoomNumber,
		RoomType:      req.RoomType,
		PricePerNight: req.PricePerNight,
		MaxGuest:      req.MaxGuest,
		Status:        req.Status,
	}, nil
}

func (d *dummyRoomService) DeleteRoom(id int) error {
	return nil
}

// ---- Tests ----

func TestCreateRoom(t *testing.T) {
	e := echo.New()
	h := handler.NewRoomHandler(&dummyRoomService{})

	body := `{"room_number":101,"room_type":"Deluxe","price_per_night":100,"max_guest":2}`
	req := httptest.NewRequest(http.MethodPost, "/rooms", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = h.CreateRoom(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "Deluxe")
}

func TestListRooms(t *testing.T) {
	e := echo.New()
	h := handler.NewRoomHandler(&dummyRoomService{})

	req := httptest.NewRequest(http.MethodGet, "/rooms", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = h.ListRooms(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Deluxe")
}

func TestGetRoomByID(t *testing.T) {
	e := echo.New()
	h := handler.NewRoomHandler(&dummyRoomService{})

	req := httptest.NewRequest(http.MethodGet, "/rooms/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	_ = h.GetRoomByID(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "101")
}
