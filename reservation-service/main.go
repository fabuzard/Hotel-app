package main

import (

	// echomw "github.com/labstack/echo/v4/middleware"

	"fmt"
	"net/http"
	"reservation-service/config"
	"reservation-service/handler"
	jwtmw "reservation-service/middleware"
	"reservation-service/model"
	"reservation-service/repository"
	"reservation-service/service"
	"reservation-service/worker"

	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
)

func main() {
	config.LoadEnv()
	db := config.DBInit()

	e := echo.New()
	r := e.Group("reservations")
	r.Use(jwtmw.JWTAuth)
	config.DB.AutoMigrate(&model.Booking{}, model.Room{})
	// room
	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo)
	roomHandler := handler.NewRoomHandler(roomService)
	// booking
	bookingRepo := repository.NewBookingRepository(db)
	bookingService := service.NewBookingService(bookingRepo)
	bookingHandler := handler.NewBookingHandler(bookingService)

	// create room routes
	r.POST("/rooms", roomHandler.CreateRoom)
	r.GET("/rooms", roomHandler.ListRooms)
	r.GET("/rooms/:id", roomHandler.GetRoomByID)
	r.PUT("/rooms/:id", roomHandler.UpdateRoom)
	r.DELETE("/rooms/:id", roomHandler.DeleteRoom)
	// booking routes
	r.POST("/bookings", bookingHandler.CreateBooking)
	r.GET("/bookings/:id", bookingHandler.GetBookingByID)
	r.PUT("/webhooks/bookings/:id", bookingHandler.WebhookUpdate)
	r.DELETE("/bookings/:id", bookingHandler.DeleteBooking)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "🚀 Server running and DB connected!")
	})
	// check in
	r.PUT("/checkin/:id", bookingHandler.Checkin)
	// check out
	r.PUT("/checkout/:id", bookingHandler.Checkout)

	// run scheduler
	c := cron.New()
	c.AddFunc("0 10 * * *", func() {
		worker.StartScheduler(db)
	})
	c.Start()
	defer c.Stop()

	fmt.Println("Connected to db")
	e.Logger.Fatal(e.Start(":8082"))
}
