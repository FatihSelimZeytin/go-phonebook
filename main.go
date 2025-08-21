package main

import (
	"go-phonebook/config"
	_ "go-phonebook/routes"
	"log/slog"
	"net/http"

	"go-phonebook/dbsql"
	_ "go-phonebook/docs"
	"go-phonebook/handlers"
	"go-phonebook/middleware"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Go Phonebook API
// @version 1.0
// @description This is the API server for the Go Phonebook app.
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load env vars
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	if err := dbsql.RunMigrations(); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		return
	}

	// Initialize DB
	db, err := dbsql.ConnectGorm()
	if err != nil {
		slog.Error("Failed to connect to DB", "error", err)
		return
	}

	// Setup Echo
	e := config.NewEchoApp()

	e.Use(middleware.CORS())

	authHandler := handlers.NewAuthHandler(db)

	e.POST("/auth/login", authHandler.Login)
	e.POST("/auth/register", authHandler.Register)

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
	if err := e.Start(":8080"); err != nil {
		slog.Error("Server failed", "error", err)
	}
}
