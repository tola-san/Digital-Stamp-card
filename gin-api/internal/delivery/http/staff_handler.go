package http

import (
	"context"
	"errors"
	"net/http"

	"gin-api/internal/domain"
	"gin-api/internal/service"

	"github.com/gin-gonic/gin"
)

const staffSessionCookie = "staff_session"

type StaffApplication interface {
	Login(context.Context, service.StaffLoginInput) (service.StaffLoginResult, error)
	Authenticate(context.Context, string) (string, error)
	Logout(context.Context, string) error
	GetStaff(context.Context, string) (domain.Staff, error)
}

type staffHandler struct {
	staff        StaffApplication
	cookieSecure bool
}

type staffLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func newStaffHandler(staff StaffApplication, cookieSecure bool) *staffHandler {
	return &staffHandler{staff: staff, cookieSecure: cookieSecure}
}

func (h *staffHandler) login(c *gin.Context) {
	var request staffLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_REQUEST", "email and password are required")
		return
	}

	result, err := h.staff.Login(c.Request.Context(), service.StaffLoginInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		writeStaffError(c, err)
		return
	}

	h.setSessionCookie(c, result.Token)
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"staff": result.Staff}})
}

func (h *staffHandler) logout(c *gin.Context) {
	token, _ := c.Cookie(staffSessionCookie)
	if err := h.staff.Logout(c.Request.Context(), token); err != nil && !errors.Is(err, domain.ErrInvalidStaffSession) {
		writeStaffError(c, err)
		return
	}

	h.clearSessionCookie(c)
	c.Status(http.StatusNoContent)
}

func (h *staffHandler) me(c *gin.Context) {
	staff, err := h.staff.GetStaff(c.Request.Context(), c.GetString(staffIDContextKey))
	if err != nil {
		writeStaffError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"staff": staff}})
}

func (h *staffHandler) setSessionCookie(c *gin.Context, token string) {
	if h.cookieSecure {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}
	c.SetCookie(
		staffSessionCookie,
		token,
		int(service.StaffSessionDuration.Seconds()),
		"/api",
		"",
		h.cookieSecure,
		true,
	)
}

func (h *staffHandler) clearSessionCookie(c *gin.Context) {
	if h.cookieSecure {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}
	c.SetCookie(staffSessionCookie, "", -1, "/api", "", h.cookieSecure, true)
}

func writeStaffError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		writeError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
	case errors.Is(err, domain.ErrInvalidStaffSession):
		writeError(c, http.StatusUnauthorized, "AUTHENTICATION_REQUIRED", "a valid staff session is required")
	case errors.Is(err, domain.ErrStaffNotFound):
		writeError(c, http.StatusNotFound, "STAFF_NOT_FOUND", err.Error())
	default:
		_ = c.Error(err)
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an internal error occurred")
	}
}
