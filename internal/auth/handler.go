package auth

import (
	"errors"
	"net/http"

	"go-first-api/internal/shared"
	"go-first-api/internal/user"

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
	{
		group.POST("/login", h.Login)
		group.GET("/me", auth, h.Me)
	}
}

func (h *Handler) Login(c *gin.Context) {
	var dto LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		shared.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.Login(c.Request.Context(), dto)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			shared.Error(c, http.StatusUnauthorized, err.Error())
			return
		}
		if err.Error() == "account is deactivated" {
			shared.Error(c, http.StatusForbidden, err.Error())
			return
		}
		shared.Error(c, http.StatusInternalServerError, "login failed")
		return
	}

	shared.OK(c, "login successful", result)
}

func (h *Handler) Me(c *gin.Context) {
	val, exists := c.Get(shared.ContextUserKey)
	if !exists {
		shared.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	shared.OK(c, "current user", val.(user.User))
}