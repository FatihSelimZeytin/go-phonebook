package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go-phonebook/models"
)

type PhoneInput struct {
	Number string `json:"number" example:"123456789"`
}

type CreateContactInput struct {
	FirstName string       `json:"firstName" example:"John"`
	Surname   string       `json:"surname" example:"Doe"`
	Phones    []PhoneInput `json:"phones"`
}

// CreateContact godoc
// @Summary      Create a new contact
// @Description  Creates a new contact with phones for the authenticated user
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Param        contact  body      CreateContactInput  true  "Contact data"
// @Success      201      {object}  models.Contact
// @Failure 	 400	  {object} utilities.BadRequestResponse
// @Failure		 500	  {object} utilities.DatabaseErrorResponse
// @Security     ApiKeyAuth
// @Router       /contacts [post]
func (h *Handler) CreateContact(c echo.Context) error {
	userID := c.Get("userID").(uint)

	var contact models.Contact
	if err := c.Bind(&contact); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
	}

	if contact.FirstName == "" || contact.Surname == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "First name and surname are required"})
	}

	// Set fields
	contact.UserID = userID
	contact.Status = true
	contact.CreatedAt = time.Now()
	contact.UpdatedAt = time.Now()

	// Save the contact first to get the ID
	if err := h.DB.Create(&contact).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create contact"})
	}

	// Assign contact ID to each phone and save
	for i := range contact.Phones {
		contact.Phones[i].ContactID = contact.ID
	}

	if len(contact.Phones) > 0 {
		if err := h.DB.Create(&contact.Phones).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save phone numbers"})
		}
	}

	return c.JSON(http.StatusCreated, contact)
}
