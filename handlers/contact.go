package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go-phonebook/models"
	"gorm.io/gorm"
)

type ContactHandler struct {
	DB *gorm.DB
}

func NewContactHandler(db *gorm.DB) *ContactHandler {
	return &ContactHandler{DB: db}
}

func (h *ContactHandler) RegisterRoutes(g *echo.Group) {
	g.GET("", h.GetAllContacts)
}

func (h *ContactHandler) GetAllContacts(c echo.Context) error {
	var contacts []models.Contact
	if err := h.DB.Find(&contacts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch contacts"})
	}
	return c.JSON(http.StatusOK, contacts)
}
