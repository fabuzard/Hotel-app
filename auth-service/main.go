// @title Shipment API
// @version 1.0
// @description This is a sample API for managing shipments.
// @host localhost:8080
// @BasePath /
// @schemes http

package main

import (
	"auth-service/config"
	"auth-service/handler"
	models "auth-service/model"
	"auth-service/repository"
	"auth-service/service"
	"fmt"

	"net/http"

	// echomw "github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	config.LoadEnv()
	db := config.DBInit()

	e := echo.New()
	config.DB.AutoMigrate(&models.User{})
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "🚀 Server running and DB connected!")
	})

	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	fmt.Println("Connected to db")
	e.Logger.Fatal(e.Start(":8080"))
}
