package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"gin-api/internal/domain"
	"gin-api/internal/service"

	"github.com/gin-gonic/gin"
)

const customerSessionCookie = "customer_session"

type CustomerApplication interface {
	Register(context.Context, service.RegisterCustomerInput) (service.RegistrationResult, error)
	GetCustomer(context.Context, string) (domain.Customer, error)
	GetCard(context.Context, string) (domain.StampCard, error)
	GetTransactions(context.Context, string, int, string) (service.TransactionPage, error)
}

type SessionApplication interface {
	Login(context.Context, string) (service.LoginResult, error)
	Authenticate(context.Context, string) (string, error)
	Logout(context.Context, string) error
}

type customerHandler struct {
	customers    CustomerApplication
	sessions     SessionApplication
	cookieSecure bool
}

type registerCustomerRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
}

type loginCustomerRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func newCustomerHandler(customers CustomerApplication, sessions SessionApplication, cookieSecure bool) *customerHandler {
	return &customerHandler{customers: customers, sessions: sessions, cookieSecure: cookieSecure}
}

func (h *customerHandler) register(c *gin.Context) {
	var request registerCustomerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_REQUEST", "name and phone are required")
		return
	}
	result, err := h.customers.Register(c.Request.Context(), service.RegisterCustomerInput{
		Name: request.Name, Phone: request.Phone,
	})
	if err != nil {
		writeCustomerError(c, err)
		return
	}
	h.setSessionCookie(c, result.Token)
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{
		"customer": result.Customer,
		"card":     result.Card,
	}})
}

func (h *customerHandler) login(c *gin.Context) {
	var request loginCustomerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_REQUEST", "phone is required")
		return
	}
	result, err := h.sessions.Login(c.Request.Context(), request.Phone)
	if err != nil {
		if errors.Is(err, domain.ErrCustomerNotFound) {
			writeError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid phone number")
			return
		}
		writeCustomerError(c, err)
		return
	}
	h.setSessionCookie(c, result.Token)
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"customer": result.Customer}})
}

func (h *customerHandler) logout(c *gin.Context) {
	token, _ := c.Cookie(customerSessionCookie)
	if err := h.sessions.Logout(c.Request.Context(), token); err != nil && !errors.Is(err, domain.ErrInvalidSession) {
		writeCustomerError(c, err)
		return
	}
	h.clearSessionCookie(c)
	c.Status(http.StatusNoContent)
}

func (h *customerHandler) me(c *gin.Context) {
	customer, err := h.customers.GetCustomer(c.Request.Context(), c.GetString(customerIDContextKey))
	if err != nil {
		writeCustomerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"customer": customer}})
}

func (h *customerHandler) card(c *gin.Context) {
	card, err := h.customers.GetCard(c.Request.Context(), c.GetString(customerIDContextKey))
	if err != nil {
		writeCustomerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"card": card}})
}

func (h *customerHandler) transactions(c *gin.Context) {
	limit := 0
	if rawLimit := c.Query("limit"); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err != nil {
			writeError(c, http.StatusUnprocessableEntity, "INVALID_PAGE_LIMIT", "limit must be between 1 and 100")
			return
		}
		limit = parsed
	}
	page, err := h.customers.GetTransactions(
		c.Request.Context(),
		c.GetString(customerIDContextKey),
		limit,
		c.Query("cursor"),
	)
	if err != nil {
		writeCustomerError(c, err)
		return
	}
	var nextCursor any
	if page.NextCursor != "" {
		nextCursor = page.NextCursor
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"transactions": page.Transactions,
		"next_cursor":  nextCursor,
	}})
}

func (h *customerHandler) setSessionCookie(c *gin.Context, token string) {
	if h.cookieSecure {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}
	c.SetCookie(
		customerSessionCookie,
		token,
		int(service.CustomerSessionDuration.Seconds()),
		"/api",
		"",
		h.cookieSecure,
		true,
	)
}

func (h *customerHandler) clearSessionCookie(c *gin.Context) {
	if h.cookieSecure {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}
	c.SetCookie(customerSessionCookie, "", -1, "/api", "", h.cookieSecure, true)
}

func writeCustomerError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrPhoneAlreadyExists):
		writeError(c, http.StatusConflict, "PHONE_ALREADY_REGISTERED", err.Error())
	case errors.Is(err, domain.ErrInvalidName):
		writeError(c, http.StatusUnprocessableEntity, "INVALID_NAME", err.Error())
	case errors.Is(err, domain.ErrInvalidPhone):
		writeError(c, http.StatusUnprocessableEntity, "INVALID_PHONE", err.Error())
	case errors.Is(err, domain.ErrInvalidCursor):
		writeError(c, http.StatusUnprocessableEntity, "INVALID_CURSOR", err.Error())
	case errors.Is(err, domain.ErrInvalidPageLimit):
		writeError(c, http.StatusUnprocessableEntity, "INVALID_PAGE_LIMIT", err.Error())
	case errors.Is(err, domain.ErrCustomerNotFound):
		writeError(c, http.StatusNotFound, "CUSTOMER_NOT_FOUND", err.Error())
	case errors.Is(err, domain.ErrInvalidSession):
		writeError(c, http.StatusUnauthorized, "AUTHENTICATION_REQUIRED", "a valid customer session is required")
	default:
		_ = c.Error(err)
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "an internal error occurred")
	}
}
