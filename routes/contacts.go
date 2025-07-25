package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) RegisterRoutes(e *echo.Group) {
	e.POST("/", h.CreateContact)
	e.GET("/", h.ListContacts)
	e.PUT("/:id", h.UpdateContact)
	e.DELETE("/:id", h.DeleteContact)
	e.GET("/search", h.SearchContacts)
}
