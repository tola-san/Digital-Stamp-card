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
	HealthChecker          HealthChecker
	FrontendURL            string
	CustomerService        CustomerApplication
	CustomerSessionService SessionApplication
	CookieSecure           bool
}

func NewRouter(deps RouterDependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), cors(deps.FrontendURL))

	api := router.Group("/api")
	api.GET("/health", healthHandler(deps.HealthChecker))
	if deps.CustomerService != nil && deps.CustomerSessionService != nil {
		customers := newCustomerHandler(deps.CustomerService, deps.CustomerSessionService, deps.CookieSecure)
		api.POST("/customers", customers.register)
		api.POST("/customer-sessions", customers.login)

		customerOnly := api.Group("")
		customerOnly.Use(customerAuth(deps.CustomerSessionService))
		customerOnly.DELETE("/customer-sessions/current", customers.logout)
		customerOnly.GET("/customers/me", customers.me)
		customerOnly.GET("/customers/me/card", customers.card)
		customerOnly.GET("/customers/me/transactions", customers.transactions)
	}

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
			"status": "Api is Running"},
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
