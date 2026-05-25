package post

import (
	"errors"
	"net/http"

	"go-first-api/internal/user"
	"go-first-api/pkg/pagination"
	"go-first-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type Middleware = gin.HandlerFunc

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *gin.Engine, authMiddleware Middleware) {
	posts := r.Group("/posts", authMiddleware)
	{
		posts.GET("", h.FindAll)
		posts.GET("/:id", h.FindOne)
		posts.GET("/user/:userId", h.FindByUserID)
		posts.POST("", h.Create)
		posts.PATCH("/:id/status", h.UpdateStatus)
		posts.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) FindAll(c *gin.Context) {
	var q pagination.Query
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	q.Normalize()
	c.JSON(http.StatusOK, h.service.FindAll(q))
}

func (h *Handler) FindOne(c *gin.Context) {
	post, err := h.service.FindOne(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.OK(c, post)
}

func (h *Handler) FindByUserID(c *gin.Context) {
	var q pagination.Query
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	q.Normalize()
	c.JSON(http.StatusOK, h.service.FindByUserID(c.Param("userId"), q))
}

func (h *Handler) Create(c *gin.Context) {
	var dto CreatePostDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	u, _ := c.Get("user")
	dto.UserID = u.(user.User).ID
	post, err := h.service.Create(dto)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Created(c, post)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	var dto UpdatePostStatusDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	post, err := h.service.UpdateStatus(c.Param("id"), dto)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, ErrInvalidStatus) {
			status = http.StatusBadRequest
		}
		response.Error(c, status, err.Error())
		return
	}
	response.OK(c, post)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Param("id")); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrNotFound) {
			status = http.StatusNotFound
		}
		response.Error(c, status, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "deleted"})
}