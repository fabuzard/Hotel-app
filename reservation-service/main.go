package main

import (

	// echomw "github.com/labstack/echo/v4/middleware"

	"fmt"
	"net/http"
	"os"
	"reservation-service/config"
	"reservation-service/handler"
	jwtmw "reservation-service/middleware"
	"reservation-service/model"
	"reservation-service/repository"
	"reservation-service/service"

	"github.com/labstack/echo/v4"
)

func main() {
	config.LoadEnv()
	db := config.DBInit()

	e := echo.New()
	r := e.Group("/api")
	r.Use(jwtmw.JWTAuth)
	config.DB.AutoMigrate(&model.Booking{}, model.Room{})
	// room
	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo)
	roomHandler := handler.NewRoomHandler(roomService)

	// create room routes
	r.POST("/rooms", roomHandler.CreateRoom)
	r.GET("/rooms", roomHandler.ListRooms)
	r.GET("/rooms/:id", roomHandler.GetRoomByID)
	// r.PUT("/rooms/:id", roomHandler.UpdateRoom)
	// r.DELETE("/rooms/:id", roomHandler.DeleteRoom)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "🚀 Server running and DB connected!")
	})

	fmt.Println("Connected to db")
	fmt.Println("JWT KEY IS = ")
	fmt.Println(os.Getenv("JWT_SECRET"))
	e.Logger.Fatal(e.Start(":8082"))
}
