package router

import (
	"gomustache/internal/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	User *handler.UserHandler
}

func New(h Handlers) *gin.Engine {
	r := gin.Default() // Logger + Recovery

	r.GET("/health", healthCheck)

	api := r.Group("/api/v1")
	h.User.Register(api)

	return r
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
