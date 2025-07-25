package config

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEchoApp() *echo.Echo {
	e := echo.New()

	e.Validator = NewValidator()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return e
}
