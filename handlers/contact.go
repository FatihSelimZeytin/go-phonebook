package handlers

import (
	"net/http"

	"go-phonebook/models"

	"github.com/labstack/echo/v4"
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

// GetAllContacts godoc
// @Summary      Get all contacts
// @Description  Get a list of all contacts for the user
// @Tags         contacts
// @Produce      json
// @Success      200  {array}   models.Contact
// @Failure		 500  {object}  map[string]interface{}
// @Router       /contacts [get]
func (h *ContactHandler) GetAllContacts(c echo.Context) error {
	var contacts []models.Contact
	if err := h.DB.Find(&contacts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch contacts"})
	}
	return c.JSON(http.StatusOK, contacts)
}

type CreateContactInput struct {
	FirstName string `json:"firstName" validate:"required,min=2" example:"John"`
	Surname   string `json:"surname" validate:"required,min=2" example:"Doe"`
	Phone     string `json:"phone" validate:"required" example:"+1234567890"`
}

// CreateContact godoc
// @Summary      Create a contact
// @Description  Create a new contact with phone number
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Param        contact  body      CreateContactInput  true  "Contact info"
// @Success      201      {object}  models.Contact
// @Failure 	 400	  {object} utilities.BadRequestResponse
// @Failure		 500	  {object} utilities.DatabaseErrorResponse
// @Router       /contacts [post]
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
		Status:    true,
		UserID:    1, // set properly from context
		Phones: []models.Phone{
			{Number: input.Phone},
		},
	}

	if err := h.DB.Create(&contact).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create contact"})
	}

	return c.JSON(http.StatusCreated, contact)
}
