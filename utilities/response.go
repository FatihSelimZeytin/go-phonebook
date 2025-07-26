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

type MessageResponse struct {
	Message string `json:"message" example:"Contact deleted"`
}

type BadRequestResponse struct {
	Error string `json:"error" example:"Bad Request"`
}

type InvalidContactIDResponse struct {
	Error string `json:"error" example:"Invalid contact ID"`
}

type UnauthorizedResponse struct {
	Error string `json:"error" example:"Unauthorized"`
}

type NotFoundResponse struct {
	Error string `json:"error" example:"Contact not found"`
}

type DatabaseErrorResponse struct {
	Error string `json:"error" example:"Database error"`
}
