package post

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
	posts := r.Group("/posts", auth)
	{
		posts.GET("", h.List)
		posts.GET("/:id", h.Get)
		posts.GET("/user/:userId", h.ListByUser)
		posts.POST("", h.Create)
		posts.PATCH("/:id/status", h.UpdateStatus)
		posts.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) List(c *gin.Context) {
	var q shared.PaginationQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		shared.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.service.FindAll(c.Request.Context(), q)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "failed to fetch posts")
		return
	}
	shared.OK(c, "posts retrieved", result)
}

func (h *Handler) Get(c *gin.Context) {
	p, err := h.service.FindByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			shared.Error(c, http.StatusNotFound, "post not found")
			return
		}
		shared.Error(c, http.StatusInternalServerError, "failed to fetch post")
		return
	}
	shared.OK(c, "post retrieved", p)
}

func (h *Handler) ListByUser(c *gin.Context) {
	var q shared.PaginationQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		shared.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.service.FindByUserID(c.Request.Context(), c.Param("userId"), q)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "failed to fetch posts")
		return
	}
	shared.OK(c, "posts retrieved", result)
}

func (h *Handler) Create(c *gin.Context) {
	caller, ok := contextUser(c)
	if !ok {
		shared.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var dto CreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		shared.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	p, err := h.service.Create(c.Request.Context(), caller.ID, dto)
	if err != nil {
		shared.Error(c, http.StatusInternalServerError, "failed to create post")
		return
	}
	shared.Created(c, "post created", p)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	caller, ok := contextUser(c)
	if !ok {
		shared.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var dto UpdateStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		shared.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	p, err := h.service.UpdateStatus(c.Request.Context(), c.Param("id"), caller.ID, dto)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			shared.Error(c, http.StatusNotFound, "post not found")
		case errors.Is(err, ErrForbidden):
			shared.Error(c, http.StatusForbidden, "you can only update your own posts")
		case errors.Is(err, ErrInvalidStatus):
			shared.Error(c, http.StatusBadRequest, "invalid status")
		default:
			shared.Error(c, http.StatusInternalServerError, "failed to update post")
		}
		return
	}
	shared.OK(c, "post updated", p)
}

func (h *Handler) Delete(c *gin.Context) {
	caller, ok := contextUser(c)
	if !ok {
		shared.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := h.service.Delete(c.Request.Context(), c.Param("id"), caller.ID); err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			shared.Error(c, http.StatusNotFound, "post not found")
		case errors.Is(err, ErrForbidden):
			shared.Error(c, http.StatusForbidden, "you can only delete your own posts")
		default:
			shared.Error(c, http.StatusInternalServerError, "failed to delete post")
		}
		return
	}
	shared.OK(c, "post deleted", nil)
}

func contextUser(c *gin.Context) (user.User, bool) {
	val, exists := c.Get(shared.ContextUserKey)
	if !exists {
		return user.User{}, false
	}
	u, ok := val.(user.User)
	return u, ok
}