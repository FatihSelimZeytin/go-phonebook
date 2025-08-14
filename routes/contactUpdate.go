package routes

import (
	"net/http"
	"strconv"
	"time"

	"go-phonebook/models"

	"github.com/labstack/echo/v4"
)

type UpdateContactInput struct {
	FirstName string        `json:"firstName" validate:"required"`
	Surname   string        `json:"surname" validate:"required"`
	Company   string        `json:"company" validate:"required"`
	Phones    []PhonesInput `json:"phones" validate:"required,dive"`
}

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

	var input UpdateContactInput

	idParam := c.Param("id")
	contactID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid contact ID"})
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
	}

	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var contact models.Contact
	if err := h.DB.Where("id = ? AND user_id = ?", contactID, userID).First(&contact).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Contact not found"})
	}

	contact.FirstName = input.FirstName
	contact.Surname = input.Surname
	contact.Company = input.Company
	contact.UpdatedAt = time.Now()

	tx := h.DB.Begin()

	if err := tx.Save(&contact).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update contact"})
	}

	// Delete old phones
	if err := tx.Where("contact_id = ?", contact.ID).Delete(&models.Phone{}).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete old phone numbers"})
	}

	// Add new phones
	for _, p := range input.Phones {
		phone := models.Phone{
			Number:    p.Number,
			ContactID: contact.ID,
		}
		if err := tx.Create(&phone).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to add phone numbers"})
		}
	}

	tx.Commit()

	// Reload contact with phones
	if err := h.DB.Preload("Phones").First(&contact, contact.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to load updated contact"})
	}

	return c.JSON(http.StatusOK, contact)
}
