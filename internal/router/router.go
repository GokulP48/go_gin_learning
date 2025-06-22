package router

import (
	"net/http"

	"github.com/GokulP48/go_gin_learning/internal/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct{}

func (s *Router) InitRouter() http.Handler {
	router := gin.Default()

	// Set CORS headers
	defaultOrigin := []string{"http://localhost:5173"}
	defaultHeaders := []string{"Accept", "Authorization", "Content-Type", "X-Custom-Header", "X-User-Id"}
	defaultMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     defaultOrigin,
		AllowMethods:     defaultMethods,
		AllowHeaders:     defaultHeaders,
		AllowCredentials: true,
	}))

	api := router.Group("/api/v1")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	api.GET("/health/db", func(c *gin.Context) {
		c.JSON(200, db.DBHealthCheck())
	})

	// api.POST("/users", s.UserHandler.CreateUser)

	return router
}
