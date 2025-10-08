// @title Shipment API
// @version 1.0
// @description This is a sample API for managing shipments.
// @host localhost:8080
// @BasePath /
// @schemes http

package main

import (
	"auth-service/config"
	models "auth-service/model"
	"fmt"

	"net/http"

	// echomw "github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	config.LoadEnv()
	config.DBInit()

	e := echo.New()
	config.DB.AutoMigrate(&models.User{})

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "🚀 Server running and DB connected!")
	})
	// e.POST("/test-error", func(c echo.Context) error {
	// 	return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
	// 		Message: "invalid body data",
	// 		Details: "missing `name` field",
	// 	})
	// })

	fmt.Println("Connected to db")
	e.Logger.Fatal(e.Start(":8080"))
}
