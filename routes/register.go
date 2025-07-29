package routes

import (
	"go-phonebook/utilities"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	PlainPassword string `json:"plainPassword"`
}

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"unique;not null" json:"username"`
	Email        string `gorm:"unique;not null" json:"email"`
	PasswordHash string `gorm:"not null"`
}

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Registers a user with username, email and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      RegisterRequest  true  "User registration data"
// @Success      201      {object}  map[string]interface{}  "User created response"
// @Failure 	 400	  {object} utilities.BadRequestResponse
// @Failure		 500	  {object} utilities.DatabaseErrorResponse
// @Router       /users/register [post]
func RegisterUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req RegisterRequest

		// Parse JSON body
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
		}

		// Validate input
		if req.Username == "" || req.Email == "" || req.PlainPassword == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "All fields are required"})
		}

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.PlainPassword), 10)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to hash password"})
		}

		// Create user
		user := User{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: string(hash),
		}

		if err := db.Create(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not create user"})
		}

		return utilities.SendCreatedUserResponse(c, user.ID, user.Username, user.Email)

	}
}
