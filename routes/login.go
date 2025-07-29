package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUser godoc
// @Summary      Login
// @Description  Authenticates a user and returns JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body LoginRequest true "Email and Password"
// @Success      200	  {object}  models.Contact
// @Failure 	 400	  {object} utilities.BadRequestResponse
// @Failure 	 401	  {object} utilities.InvalidContactIDResponse
// @Failure		 500	  {object} utilities.DatabaseErrorResponse
// @Router       /auth/login [post]
func LoginUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest

		// Bind and validate input
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
		}
		if req.Email == "" || req.Password == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email and password are required"})
		}

		// Find user
		var user User
		if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
		}

		// Check password
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
		}

		// Create JWT token
		claims := jwt.MapClaims{
			"user_id": user.ID,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secret := os.Getenv("JWT_SECRET")
		signedToken, err := token.SignedString([]byte(secret))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to sign token"})
		}

		return c.JSON(http.StatusOK, echo.Map{"token": signedToken})
	}
}
