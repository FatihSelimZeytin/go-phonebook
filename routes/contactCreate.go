package routes

import (
	"net/http"
	"time"

	"go-phonebook/models"

	"github.com/labstack/echo/v4"
)

type PhonesInput struct {
	Number string `json:"number" validate:"required"`
}

type CreateContactInput struct {
	FirstName string        `json:"firstName" validate:"required"`
	Surname   string        `json:"surname" validate:"required"`
	Company   string        `json:"company" validate:"required"`
	Phones    []PhonesInput `json:"phones" validate:"required,dive"`
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

	var input CreateContactInput

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
	}

	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if input.FirstName == "" || input.Surname == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "First name and surname are required"})
	}

	// Map input to your model.Contact
	contact := models.Contact{
		FirstName: input.FirstName,
		Surname:   input.Surname,
		Company:   input.Company,
		UserID:    userID,
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Map phones input to models.Phones slice
	for _, p := range input.Phones {
		contact.Phones = append(contact.Phones, models.Phone{
			Number: p.Number,
		})
	}

	// Save contact (this should save phones too if association is setups correctly)
	if err := h.DB.Create(&contact).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create contact"})
	}

	return c.JSON(http.StatusCreated, contact)
}
