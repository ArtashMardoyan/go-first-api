package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"go-first-api/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Handler struct {
	userService *user.Service
}

func NewHandler(userService *user.Service) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.GET("/me", JWTMiddleware(h.userService), h.Me)
	}
}

func (h *Handler) Login(c *gin.Context) {
	var dto LoginDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	u, err := h.userService.FindByEmail(dto.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": ErrInvalidCredentials.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(dto.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": ErrInvalidCredentials.Error()})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{AccessToken: signed})
}

func (h *Handler) Me(c *gin.Context) {
	u, _ := c.Get(ContextUser)
	c.JSON(http.StatusOK, u)
}