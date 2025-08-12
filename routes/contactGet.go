package routes

import (
	"go-phonebook/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetContact godoc
// @Summary      Get a contact by ID
// @Description  Get contact details by contact ID for the authenticated user
// @Tags         contacts
// @Produce      json
// @Success      200   {object} models.Contact
// @Failure      404   {object} map[string]string
// @Failure      500   {object} map[string]string
// @Security     ApiKeyAuth
// @Router       /contacts/{id} [get]
func (h *Handler) GetContact(c echo.Context) error {
	userID := c.Get("userID").(uint)
	id := c.Param("id")

	var contact models.Contact
	// Find contact with given id and user_id and status=true
	if err := h.DB.Preload("Phones").
		Where("id = ? AND user_id = ? AND status = ?", id, userID, true).
		First(&contact).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Contact not found"})
	}

	return c.JSON(http.StatusOK, contact)
}
