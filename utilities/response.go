package utilities

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func SendCreatedUserResponse(c echo.Context, id uint, username, email string) error {
	return c.JSON(http.StatusCreated, echo.Map{
		"id":       id,
		"username": username,
		"email":    email,
	})
}
