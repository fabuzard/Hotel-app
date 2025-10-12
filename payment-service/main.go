package main

import (

	// echomw "github.com/labstack/echo/v4/middleware"

	"fmt"
	"net/http"
	"os"
	"payment-service/config"
	"payment-service/handler"
	jwtmw "payment-service/middleware"
	"payment-service/model"
	"payment-service/repository"
	"payment-service/service"

	"github.com/labstack/echo/v4"
)

func main() {
	config.LoadEnv()
	db := config.DBInit()

	e := echo.New()
	r := e.Group("payments")
	r.Use(jwtmw.JWTAuth)
	config.DB.AutoMigrate(&model.Payment{})
	// payment
	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "🚀 Server running and DB connected!")
	})

	r.GET("/test", func(c echo.Context) error {
		userID := c.Get("user_id")
		username := c.Get("username")
		return c.JSON(200, map[string]interface{}{
			"user_id":  userID,
			"username": username,
			"message":  "You are in a restricted area",
		})
	})

	// create payment routes
	r.POST("/pay", paymentHandler.CreatePayment)
	r.POST("/webhook", paymentHandler.SimulateWebhook)
	fmt.Println("Connected to db")
	fmt.Println("JWT KEY IS = ")
	fmt.Println(os.Getenv("JWT_SECRET"))
	e.Logger.Fatal(e.Start(":8083"))
}
