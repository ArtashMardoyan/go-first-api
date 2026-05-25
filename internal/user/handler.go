package user

import (
	"errors"
	"net/http"

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
	var q pagination.Query
	if err := c.ShouldBindQuery(&q); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	q.Normalize()
	c.JSON(http.StatusOK, h.service.FindAll(q))
}

func (h *Handler) FindOne(c *gin.Context) {
	user, err := h.service.FindOne(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	response.OK(c, user)
}

func (h *Handler) Create(c *gin.Context) {
	var dto CreateUserDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.service.Create(dto)
	if err != nil {
		response.Error(c, statusFromErr(err), err.Error())
		return
	}
	response.Created(c, user)
}

func contextUser(c *gin.Context) (User, bool) {
	u, ok := c.Get("user")
	if !ok {
		return User{}, false
	}
	user, ok := u.(User)
	return user, ok
}

func (h *Handler) Update(c *gin.Context) {
	cu, ok := contextUser(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var dto UpdateUserDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.service.Update(cu.ID, dto)
	if err != nil {
		response.Error(c, statusFromErr(err), err.Error())
		return
	}
	response.OK(c, user)
}

func (h *Handler) Delete(c *gin.Context) {
	cu, ok := contextUser(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := h.service.Delete(cu.ID); err != nil {
		response.Error(c, statusFromErr(err), err.Error())
		return
	}
	response.OK(c, gin.H{"message": "deleted"})
}

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