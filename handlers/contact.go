package handlers

import (
	"go-phonebook/routes"
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
	g.GET("", h.ListContacts)
	g.POST("", h.CreateContact)
	g.GET("/:id", h.GetContact)
	g.PUT("/:id", h.UpdateContact)
	g.DELETE("/:id", h.DeleteContact)
	g.GET("/search", h.SearchContacts)
}

func (h *ContactHandler) GetAllContacts(c echo.Context) error {
	var contacts []models.Contact
	if err := h.DB.Find(&contacts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch contacts"})
	}
	return c.JSON(http.StatusOK, contacts)
}

func (h *ContactHandler) CreateContact(c echo.Context) error {
	routesHandler := &routes.Handler{DB: h.DB}
	return routesHandler.CreateContact(c)
}

func (h *ContactHandler) DeleteContact(c echo.Context) error {
	routesHandler := &routes.Handler{DB: h.DB}
	return routesHandler.DeleteContact(c)
}

func (h *ContactHandler) ListContacts(c echo.Context) error {
	routesHandler := &routes.Handler{DB: h.DB}
	return routesHandler.ListContacts(c)
}

func (h *ContactHandler) GetContact(c echo.Context) error {
	routesHandler := &routes.Handler{DB: h.DB}
	return routesHandler.GetContact(c)
}

func (h *ContactHandler) SearchContacts(c echo.Context) error {
	routesHandler := &routes.Handler{DB: h.DB}
	return routesHandler.SearchContacts(c)
}

func (h *ContactHandler) UpdateContact(c echo.Context) error {
	routesHandler := &routes.Handler{DB: h.DB}
	return routesHandler.UpdateContact(c)
}
