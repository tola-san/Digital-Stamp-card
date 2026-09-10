package http

import (
	"errors"
	"net/http"

	"gin-api/internal/domain"

	"github.com/gin-gonic/gin"
)

const customerIDContextKey = "customer_id"

func customerAuth(sessions SessionApplication) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(customerSessionCookie)
		if err != nil || token == "" {
			writeError(c, http.StatusUnauthorized, "AUTHENTICATION_REQUIRED", "a valid customer session is required")
			return
		}
		customerID, err := sessions.Authenticate(c.Request.Context(), token)
		if err != nil {
			if !errors.Is(err, domain.ErrInvalidSession) {
				_ = c.Error(err)
			}
			writeError(c, http.StatusUnauthorized, "AUTHENTICATION_REQUIRED", "a valid customer session is required")
			return
		}
		c.Set(customerIDContextKey, customerID)
		c.Next()
	}
}
