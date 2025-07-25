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
	g.POST("", h.CreateContact)
}

func (h *ContactHandler) GetAllContacts(c echo.Context) error {
	var contacts []models.Contact
	if err := h.DB.Find(&contacts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch contacts"})
	}
	return c.JSON(http.StatusOK, contacts)
}

type CreateContactInput struct {
	FirstName string `json:"firstName" validate:"required,min=2"`
	Surname   string `json:"surname" validate:"required,min=2"`
	Phone     string `json:"phone" validate:"required"`
}

func (h *ContactHandler) CreateContact(c echo.Context) error {
	var input CreateContactInput

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	if err := c.Validate(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	contact := models.Contact{
		FirstName: input.FirstName,
		Surname:   input.Surname,
		Status:    true, // assuming default active
		UserID:    1,    // set to a proper user from context later
		Phones: []models.Phone{
			{Number: input.Phone},
		},
	}

	if err := h.DB.Create(&contact).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create contact"})
	}

	return c.JSON(http.StatusCreated, contact)
}
