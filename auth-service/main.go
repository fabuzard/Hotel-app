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

	jwtmw "auth-service/middleware"
	"net/http"

	// echomw "github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	config.LoadEnv()
	db := config.DBInit()

	e := echo.New()
	config.DB.AutoMigrate(&models.Users{})
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authHandler := handler.NewAuthHandler(authService)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "🚀 Server running and DB connected!")
	})

	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	// test protected route
	r := e.Group("/restricted")
	r.Use(jwtmw.JWTAuth)
	r.GET("", func(c echo.Context) error {
		userID := c.Get("user_id")
		username := c.Get("username")
		return c.JSON(200, map[string]interface{}{
			"user_id":  userID,
			"username": username,
			"message":  "You are in a restricted area",
		})
	})
	fmt.Println("Connected to db")
	e.Logger.Fatal(e.Start(":8081"))
}
