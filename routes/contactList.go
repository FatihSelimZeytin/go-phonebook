package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go-phonebook/models"
)

func (h *Handler) ListContacts(c echo.Context) error {
	userID := c.Get("userID").(uint)

	var contacts []models.Contact
	if err := h.DB.Preload("Phones").
		Where("user_id = ? AND status = ?", userID, true).
		Find(&contacts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, contacts)
}
