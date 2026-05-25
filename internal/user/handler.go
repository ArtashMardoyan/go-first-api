package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Middleware = gin.HandlerFunc

// аналог @Controller() в NestJS
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine, authMiddleware Middleware) {
	users := r.Group("/users")
	{
		users.GET("", authMiddleware, h.FindAll)
		users.GET("/:id", authMiddleware, h.FindOne)
		users.POST("", h.Create)
		users.PATCH("", authMiddleware, h.Update)
		users.DELETE("", authMiddleware, h.Delete)
	}
}

func (h *Handler) FindAll(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.FindAll())
}

func (h *Handler) FindOne(c *gin.Context) {
	user, err := h.service.FindOne(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) Create(c *gin.Context) {
	var dto CreateUserDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err := h.service.Create(dto)
	if err != nil {
		c.JSON(statusFromErr(err), gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func contextUser(c *gin.Context) User {
	u, _ := c.Get("user")
	return u.(User)
}

func (h *Handler) Update(c *gin.Context) {
	var dto UpdateUserDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err := h.service.Update(contextUser(c).ID, dto)
	if err != nil {
		c.JSON(statusFromErr(err), gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(contextUser(c).ID); err != nil {
		c.JSON(statusFromErr(err), gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// statusFromErr — маппинг ошибок на HTTP-статусы, аналог ExceptionFilter в NestJS
func statusFromErr(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrEmailTaken):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
