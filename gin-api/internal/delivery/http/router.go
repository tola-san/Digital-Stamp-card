package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthChecker interface {
	Ping(context.Context) error
}

type RouterDependencies struct {
	HealthChecker HealthChecker
	FrontendURL   string
}

func NewRouter(deps RouterDependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), cors(deps.FrontendURL))

	api := router.Group("/api")
	api.GET("/health", healthHandler(deps.HealthChecker))

	router.NoRoute(func(c *gin.Context) {
		writeError(c, http.StatusNotFound, "NOT_FOUND", "route not found")
	})

	return router
}

func healthHandler(checker HealthChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		if checker == nil {
			writeError(c, http.StatusServiceUnavailable, "DATABASE_UNAVAILABLE", "database health check is not configured")
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := checker.Ping(ctx); err != nil {
			writeError(c, http.StatusServiceUnavailable, "DATABASE_UNAVAILABLE", "database is unavailable")
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": gin.H{
			"status": "api is running"},
		},
		)
	}
}

func cors(frontendURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if frontendURL != "" && origin == frontendURL {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Vary", "Origin")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
