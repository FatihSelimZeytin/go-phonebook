package routes

import (
	"errors"
	"go-phonebook/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	_ "gorm.io/gorm/clause"
)

type Handler struct {
	DB *gorm.DB
}

func (h *Handler) RegisterRoutes(e *echo.Group) {
	e.POST("/", h.CreateContact)
	e.GET("/", h.ListContacts)
	e.PUT("/", h.UpdateContact)
	e.DELETE("/", h.DeleteContact)
	e.GET("/search", h.SearchContacts)
}

// CreateContact creates a contact with multiple phones
func (h *Handler) CreateContact(c echo.Context) error {
	userID := c.Get("userID").(uint)

	var req struct {
		FirstName string   `json:"firstName" validate:"required"`
		Surname   string   `json:"surname" validate:"required"`
		Company   string   `json:"company"`
		Phones    []string `json:"phones"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	contact := models.Contact{
		FirstName: req.FirstName,
		Surname:   req.Surname,
		Company:   req.Company,
		UserID:    userID,
		Status:    true,
	}

	if err := h.DB.Create(&contact).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Add phones if any
	for _, number := range req.Phones {
		phone := models.Phone{
			Number:    number,
			ContactID: contact.ID,
		}
		if err := h.DB.Create(&phone).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
	}

	// Load contact with phones
	h.DB.Preload("Phones").First(&contact, contact.ID)

	return c.JSON(http.StatusCreated, contact)
}

// ListContacts returns all active contacts for current user
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

// UpdateContact updates a contact found by firstName and surname (case-insensitive)
func (h *Handler) UpdateContact(c echo.Context) error {
	userID := c.Get("userID").(uint)

	var req struct {
		SearchFirstName string   `json:"searchFirstName" validate:"required"`
		SearchSurname   string   `json:"searchSurname" validate:"required"`
		FirstName       string   `json:"firstName"`
		Surname         string   `json:"surname"`
		Company         string   `json:"company"`
		Phones          []string `json:"phones"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if req.SearchFirstName == "" || req.SearchSurname == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "searchFirstName and searchSurname required"})
	}

	var contact models.Contact
	if err := h.DB.Preload("Phones").Where("user_id = ? AND LOWER(first_name) = ? AND LOWER(surname) = ?", userID, strings.ToLower(req.SearchFirstName), strings.ToLower(req.SearchSurname)).First(&contact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Contact not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	updates := map[string]interface{}{}
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.Surname != "" {
		updates["surname"] = req.Surname
	}
	if req.Company != "" {
		updates["company"] = req.Company
	}
	if len(updates) > 0 {
		if err := h.DB.Model(&contact).Updates(updates).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	if req.Phones != nil {
		// Delete old phones
		if err := h.DB.Where("contact_id = ?", contact.ID).Delete(&models.Phone{}).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		// Add new phones
		for _, number := range req.Phones {
			phone := models.Phone{
				Number:    number,
				ContactID: contact.ID,
			}
			if err := h.DB.Create(&phone).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}
	}

	h.DB.Preload("Phones").First(&contact, contact.ID)
	return c.JSON(http.StatusOK, contact)
}

// DeleteContact sets status = false (soft delete)
func (h *Handler) DeleteContact(c echo.Context) error {
	userID := c.Get("userID").(uint)
	firstName := c.QueryParam("firstName")
	surname := c.QueryParam("surname")

	if firstName == "" || surname == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "firstName and surname required"})
	}

	var contact models.Contact
	if err := h.DB.Where("user_id = ? AND status = ? AND LOWER(first_name) = ? AND LOWER(surname) = ?", userID, true, strings.ToLower(firstName), strings.ToLower(surname)).First(&contact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Contact not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	contact.Status = false
	if err := h.DB.Save(&contact).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Contact disabled (soft deleted)"})
}

// SearchContacts searches contacts by query string q
func (h *Handler) SearchContacts(c echo.Context) error {
	userID := c.Get("userID").(uint)
	q := c.QueryParam("q")

	if q == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Search query is required"})
	}

	query := strings.ToLower(strings.TrimSpace(q))
	queryParts := strings.Fields(query)

	dbQuery := h.DB.Preload("Phones").Where("user_id = ? AND status = ?", userID, true).Where(func(db *gorm.DB) *gorm.DB {
		orCond := db
		for _, part := range queryParts {
			like := "%" + part + "%"
			orCond = orCond.Or(
				"LOWER(first_name) LIKE ? OR LOWER(surname) LIKE ? OR LOWER(company) LIKE ? OR EXISTS (SELECT 1 FROM phones WHERE contact_id = contacts.id AND number LIKE ?)",
				like, like, like, like,
			)
		}
		return orCond
	})

	var contacts []models.Contact
	if err := dbQuery.Find(&contacts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, contacts)
}
