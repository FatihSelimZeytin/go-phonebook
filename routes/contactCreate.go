package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go-phonebook/models"
)

func (h *Handler) CreateContact(c echo.Context) error {
	userID := c.Get("userID").(uint)

	var contact models.Contact
	if err := c.Bind(&contact); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
	}

	// Required field checks (optional, you can enforce this on model level too)
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
