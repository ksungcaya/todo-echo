package main

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/ksungcaya/todo-echo/configs"
	"github.com/ksungcaya/todo-echo/controllers"
	"github.com/ksungcaya/todo-echo/database"
	"github.com/ksungcaya/todo-echo/repositories"
	"github.com/ksungcaya/todo-echo/router"
	"github.com/labstack/echo/v4"
)

func main() {
	config := configs.New()
	db := database.New(&config.Database, config.IsProd())
	if conn, ok := db.DB(); ok != nil {
		defer conn.Close()
	}

	database.AutoMigrate(db)
	// database.Refresh(db)

	userRepo := repositories.NewUserRepository(db)
	authController := controllers.NewAuth(userRepo)

	r := router.New()
	r.GET("/", hello)
	r.POST("/auth/register", authController.Register)

	// Start server
	// r.Logger.Fatal(r.Start(fmt.Sprintf(":%d", config.Port)))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
