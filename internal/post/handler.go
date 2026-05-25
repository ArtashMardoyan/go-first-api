package post

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	posts := r.Group("/posts")
	{
		posts.GET("", h.FindAll)
		posts.GET("/:id", h.FindOne)
		posts.GET("/user/:userId", h.FindByUserID)
		posts.POST("", h.Create)
		posts.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) FindAll(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.FindAll())
}

func (h *Handler) FindOne(c *gin.Context) {
	post, err := h.service.FindOne(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}

func (h *Handler) FindByUserID(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.FindByUserID(c.Param("userId")))
}

func (h *Handler) Create(c *gin.Context) {
	var dto CreatePostDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	post, err := h.service.Create(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, post)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}