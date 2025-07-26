package main

import (
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go-phonebook/config"
	"go-phonebook/dbsql"
	_ "go-phonebook/docs"
	"go-phonebook/handlers"
	"go-phonebook/middleware"
)

// @title Go Phonebook API
// @version 1.0
// @description This is the API server for the Go Phonebook app.
// @host localhost:8099
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load env vars
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	// Initialize DB
	db, err := dbsql.ConnectDB()
	if err != nil {
		slog.Error("Failed to connect to DB", "error", err)
		return
	}

	// Setup Echo
	e := config.NewEchoApp()

	// Root route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World! Test")
	})

	// JWT Protected Group
	api := e.Group("/contacts")
	api.Use(middleware.JWTMiddleware)

	contactHandler := handlers.NewContactHandler(db)
	contactHandler.RegisterRoutes(api)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// Start server
	if err := e.Start(":8091"); err != nil {
		slog.Error("Server failed", "error", err)
	}
}
