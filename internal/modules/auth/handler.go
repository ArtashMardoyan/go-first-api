package auth

import (
	"errors"
	"net/http"

	"go-first-api/internal/modules/user"
	"go-first-api/internal/shared"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine, auth gin.HandlerFunc) {
	group := r.Group("/auth")
	group.POST("/login", h.Login)
	group.GET("/me", auth, h.Me)
}

func (h *Handler) Login(c *gin.Context) {
	var dto LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		shared.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Login(c.Request.Context(), dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			shared.Error(c, http.StatusUnauthorized, "invalid credentials")
		case errors.Is(err, ErrAccountDeactivated):
			shared.Error(c, http.StatusForbidden, "account is deactivated")
		default:
			shared.Error(c, http.StatusInternalServerError, "login failed")
		}
		return
	}

	shared.OK(c, "login successful", result)
}

func (*Handler) Me(c *gin.Context) {
	val, exists := c.Get(shared.ContextUserKey)
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	u, ok := val.(user.User)
	if !ok {
		shared.Error(c, http.StatusInternalServerError, "internal error")
		return
	}
	shared.OK(c, "current user", u)
}
