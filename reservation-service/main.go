package main

import (

	// echomw "github.com/labstack/echo/v4/middleware"

	"fmt"
	"net/http"
	"reservation-service/config"
	"reservation-service/handler"
	"reservation-service/model"
	"reservation-service/repository"
	"reservation-service/service"

	"github.com/labstack/echo/v4"
)

func main() {
	config.LoadEnv()
	db := config.DBInit()

	e := echo.New()
	config.DB.AutoMigrate(&model.Booking{}, model.Room{})
	// room
	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo)
	roomHandler := handler.NewRoomHandler(roomService)

	// create room routes
	e.POST("/rooms", roomHandler.CreateRoom)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "🚀 Server running and DB connected!")
	})

	fmt.Println("Connected to db")
	e.Logger.Fatal(e.Start(":8082"))
}
