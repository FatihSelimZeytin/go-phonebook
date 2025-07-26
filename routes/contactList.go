package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go-phonebook/models"
)

// ListContacts godoc
// @Summary      List all active contacts
// @Description  Get all contacts with status=true for the authenticated user, including phone numbers
// @Tags         contacts
// @Produce      json
// @Success      200  	  {array}   models.Contact
// @Failure		 500	  {object} utilities.DatabaseErrorResponse
// @Security     ApiKeyAuth
// @Router       /contacts [get]
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
