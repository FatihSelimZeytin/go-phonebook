package routes

import (
	"net/http"
	"strings"

	"go-phonebook/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SearchContacts godoc
// @Summary      Search contacts
// @Description  Search user's contacts by first name, surname, company, or phone number (case-insensitive, partial match)
// @Tags         contacts
// @Security     ApiKeyAuth
// @Param        q   query     string  true  "Search query string"
// @Success      200 {array}   models.Contact
// @Failure 	 400	  {object} utilities.BadRequestResponse
// @Failure		 500	  {object} utilities.DatabaseErrorResponse
// @Router       /contacts/search [get]
func (h *Handler) SearchContacts(c echo.Context) error {
	userID := c.Get("userID").(uint)
	q := c.QueryParam("q")

	if q == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Search query is required"})
	}

	// Prepare search terms: split by space and lowercase
	query := strings.ToLower(strings.TrimSpace(q))
	queryParts := strings.Fields(query)

	// Build a complex query with OR conditions for each part matching any of the searchable fields
	dbQuery := h.DB.Preload("Phones").
		Where("user_id = ? AND status = ?", userID, true).
		Where(func(db *gorm.DB) *gorm.DB {
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
