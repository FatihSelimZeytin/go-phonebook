package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go-phonebook/models"
)

// UpdateContact godoc
// @Summary      Update a contact
// @Description  Update contact details including phones for the authenticated user
// @Tags         contacts
// @Security     ApiKeyAuth
// @Param        id     path      int            true  "Contact ID"
// @Param        contact  body     models.Contact true  "Contact data to update"
// @Success      200	  {object}  models.Contact
// @Failure 	 400	  {object} utilities.BadRequestResponse
// @Failure		 404	  {object} utilities.NotFoundResponse
// @Failure		 500	  {object} utilities.DatabaseErrorResponse
// @Router       /contacts/{id} [put]
func (h *Handler) UpdateContact(c echo.Context) error {
	userID := c.Get("userID").(uint)

	// Parse contact ID from route
	idParam := c.Param("id")
	contactID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid contact ID"})
	}

	var updatedData models.Contact
	if err := c.Bind(&updatedData); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request data"})
	}

	var contact models.Contact
	if err := h.DB.Where("id = ? AND user_id = ?", contactID, userID).First(&contact).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Contact not found"})
	}

	// Update fields
	contact.FirstName = updatedData.FirstName
	contact.Surname = updatedData.Surname
	contact.Company = updatedData.Company
	contact.UpdatedAt = time.Now()

	if err := h.DB.Save(&contact).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update contact"})
	}

	// Remove old phone numbers
	h.DB.Where("contact_id = ?", contact.ID).Delete(&models.Phone{})

	// Add new phone numbers
	for _, p := range updatedData.Phones {
		p.ContactID = contact.ID
		h.DB.Create(&p)
	}

	return c.JSON(http.StatusOK, contact)
}
