package routes

import (
	"errors"
	"net/http"
	"strconv"

	"go-phonebook/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// DeleteContact godoc
// @Summary      Delete a contact
// @Description  Delete a contact by ID for the authenticated user
// @Tags         contacts
// @Produce      json
// @Param        id   path      int  true  "Contact ID"
// @Success      200  	  {object} utilities.MessageResponse
// @Failure 	 401	  {object} utilities.InvalidContactIDResponse
// @Failure 	 402	  {object} utilities.UnauthorizedResponse
// @Failure		 404	  {object} utilities.NotFoundResponse
// @Failure		 500	  {object} utilities.DatabaseErrorResponse
// @Security     ApiKeyAuth
// @Router       /contacts/{id} [delete]
func (h *Handler) DeleteContact(c echo.Context) error {
	// Parse contact ID from URL
	contactID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid contact ID"})
	}

	// Get user ID from context (JWT middleware should set this)
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	// Find contact and ensure it belongs to the user
	var contact models.Contact
	if err := h.DB.Where("id = ? AND user_id = ?", contactID, userID).First(&contact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Contact not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
	}

	// Delete contact (also deletes phones if you have ON DELETE CASCADE set in DB or use manual deletion)
	contact.Status = false
	if err := h.DB.Save(&contact).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to soft delete contact"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Contact deleted"})
}
