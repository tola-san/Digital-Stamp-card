package http

import (
	"errors"
	"net/http"

	"gin-api/internal/domain"

	"github.com/gin-gonic/gin"
)

const staffIDContextKey = "staff_id"

func staffAuth(sessions StaffApplication) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(staffSessionCookie)
		if err != nil || token == "" {
			writeError(c, http.StatusUnauthorized, "AUTHENTICATION_REQUIRED", "a valid staff session is required")
			return
		}

		staffID, err := sessions.Authenticate(c.Request.Context(), token)
		if err != nil {
			if !errors.Is(err, domain.ErrInvalidStaffSession) {
				_ = c.Error(err)
			}
			writeError(c, http.StatusUnauthorized, "AUTHENTICATION_REQUIRED", "a valid staff session is required")
			return
		}

		c.Set(staffIDContextKey, staffID)
		c.Next()
	}
}
