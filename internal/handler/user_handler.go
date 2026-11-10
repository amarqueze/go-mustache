package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"gomustache/internal/model"
	"gomustache/internal/service"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	ListUsers(ctx context.Context) ([]model.User, error)
	GetUser(ctx context.Context, id int) (*model.User, error)
	GetUserWithPosts(ctx context.Context, id int) (*model.UserWithPosts, error)
	CreateUser(ctx context.Context, in model.UserInput) error
	UpdateUser(ctx context.Context, id int, in model.UserInput) error
}

var _ UserService = (*service.UserService)(nil)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.GET("", h.ListUsers)
	users.GET("/:id", h.GetUser)
	users.GET("/:id/posts", h.GetUserPosts)
	users.POST("", h.CreateUser)
	users.PUT("/:id", h.UpdateUser)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.svc.ListUsers(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}

	out := make([]model.UserSummary, 0, len(users))
	for _, u := range users {
		out = append(out, model.UserSummary{ID: u.ID, Name: u.Name, Email: u.Email})
	}
	c.JSON(http.StatusOK, out)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	u, err := h.svc.GetUser(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) GetUserPosts(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	res, err := h.svc.GetUserWithPosts(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var in model.UserInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.CreateUser(c.Request.Context(), in); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "ok"})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var in model.UserInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateUser(c.Request.Context(), id, in); err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func parseID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return 0, false
	}
	return id, true
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	case errors.Is(err, context.DeadlineExceeded):
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "upstream timeout"})
	default:
		log.Printf("internal error on %s %s: %v", c.Request.Method, c.FullPath(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}
