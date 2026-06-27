package health

import (
	"go-first-api/internal/shared"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", h.Ping)
}

func (*Handler) Ping(c *gin.Context) {
	shared.OK(c, "pong")
}
