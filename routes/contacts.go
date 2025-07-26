package routes

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

// RegisterRoutes godoc
// @Summary      Register contact routes
// @Description  Register CRUD routes for contacts: create, list, update, delete, search
// @Tags         contacts
// @Security     ApiKeyAuth
func (h *Handler) RegisterRoutes(e *echo.Group) {
	e.POST("/", h.CreateContact)
	e.GET("/", h.ListContacts)
	e.PUT("/:id", h.UpdateContact)
	e.DELETE("/:id", h.DeleteContact)
	e.GET("/search", h.SearchContacts)
}
