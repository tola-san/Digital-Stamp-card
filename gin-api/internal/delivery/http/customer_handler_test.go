package http_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httpdelivery "gin-api/internal/delivery/http"
	"gin-api/internal/domain"
	"gin-api/internal/service"
)

type customerApplicationStub struct {
	registerResult service.RegistrationResult
	registerErr    error
	customer       domain.Customer
	card           domain.StampCard
	requestedID    string
}

func (s *customerApplicationStub) Register(
	context.Context,
	service.RegisterCustomerInput,
) (service.RegistrationResult, error) {
	return s.registerResult, s.registerErr
}

func (s *customerApplicationStub) GetCustomer(_ context.Context, customerID string) (domain.Customer, error) {
	s.requestedID = customerID
	return s.customer, nil
}

func (s *customerApplicationStub) GetCard(_ context.Context, customerID string) (domain.StampCard, error) {
	s.requestedID = customerID
	return s.card, nil
}

func (s *customerApplicationStub) GetTransactions(
	context.Context,
	string,
	int,
	string,
) (service.TransactionPage, error) {
	return service.TransactionPage{Transactions: []domain.StampTransaction{}}, nil
}

type sessionApplicationStub struct {
	loginResult           service.LoginResult
	loginErr              error
	authenticatedID       string
	authenticateErr       error
	authenticatedRawToken string
	logoutRawToken        string
}

func (s *sessionApplicationStub) Login(context.Context, string) (service.LoginResult, error) {
	return s.loginResult, s.loginErr
}

func (s *sessionApplicationStub) Authenticate(_ context.Context, rawToken string) (string, error) {
	s.authenticatedRawToken = rawToken
	return s.authenticatedID, s.authenticateErr
}

func (s *sessionApplicationStub) Logout(_ context.Context, rawToken string) error {
	s.logoutRawToken = rawToken
	return nil
}

func newCustomerRouter(customers *customerApplicationStub, sessions *sessionApplicationStub) http.Handler {
	return httpdelivery.NewRouter(httpdelivery.RouterDependencies{
		HealthChecker:          healthChecker{},
		CustomerService:        customers,
		CustomerSessionService: sessions,
	})
}

func TestRegisterCustomerCreatesSessionCookie(t *testing.T) {
	customers := &customerApplicationStub{registerResult: service.RegistrationResult{
		Customer: domain.Customer{ID: "customer-1", Name: "Tola San", Phone: "+85512345678"},
		Card:     domain.StampCard{ID: "card-1", CustomerID: "customer-1", RequiredStamps: 10},
		Token:    "raw-session-token",
	}}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewBufferString(
		`{"name":"Tola San","phone":"+85512345678"}`,
	))
	request.Header.Set("Content-Type", "application/json")
	newCustomerRouter(customers, &sessionApplicationStub{}).ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d: %s", http.StatusCreated, response.Code, response.Body.String())
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != "customer_session" || cookies[0].Value != "raw-session-token" {
		t.Fatalf("unexpected registration cookie: %+v", cookies)
	}
	if !cookies[0].HttpOnly || cookies[0].Path != "/api" {
		t.Fatalf("registration cookie is not properly restricted: %+v", cookies[0])
	}
}

func TestRegisterCustomerMapsDuplicatePhoneToConflict(t *testing.T) {
	customers := &customerApplicationStub{registerErr: domain.ErrPhoneAlreadyExists}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/customers", bytes.NewBufferString(
		`{"name":"Tola San","phone":"+85512345678"}`,
	))
	request.Header.Set("Content-Type", "application/json")
	newCustomerRouter(customers, &sessionApplicationStub{}).ServeHTTP(response, request)

	if response.Code != http.StatusConflict || !strings.Contains(response.Body.String(), "PHONE_ALREADY_REGISTERED") {
		t.Fatalf("unexpected duplicate response: %d %s", response.Code, response.Body.String())
	}
}

func TestProtectedCustomerCardUsesAuthenticatedCustomer(t *testing.T) {

	customers := &customerApplicationStub{card: domain.StampCard{
		ID: "card-1", CustomerID: "customer-1", StampCount: 4, RequiredStamps: 10,
	}}
	sessions := &sessionApplicationStub{authenticatedID: "customer-1"}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/customers/me/card", nil)
	request.AddCookie(&http.Cookie{Name: "customer_session", Value: "raw-session-token"})
	newCustomerRouter(customers, sessions).ServeHTTP(response, request)

	if response.Code != http.StatusOK || customers.requestedID != "customer-1" {
		t.Fatalf("unexpected card response: %d %s", response.Code, response.Body.String())
	}
	if sessions.authenticatedRawToken != "raw-session-token" {
		t.Fatalf("middleware received %q", sessions.authenticatedRawToken)
	}
}

func TestProtectedCustomerRouteRejectsMissingSession(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/customers/me", nil)
	newCustomerRouter(&customerApplicationStub{}, &sessionApplicationStub{}).ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized || !strings.Contains(response.Body.String(), "AUTHENTICATION_REQUIRED") {
		t.Fatalf("unexpected unauthorized response: %d %s", response.Code, response.Body.String())
	}
}

func TestPhoneLoginDoesNotRevealUnknownCustomer(t *testing.T) {
	sessions := &sessionApplicationStub{loginErr: domain.ErrCustomerNotFound}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/customer-sessions", bytes.NewBufferString(
		`{"phone":"+85512345678"}`,
	))
	request.Header.Set("Content-Type", "application/json")
	newCustomerRouter(&customerApplicationStub{}, sessions).ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized || !strings.Contains(response.Body.String(), "INVALID_CREDENTIALS") {
		t.Fatalf("unexpected login response: %d %s", response.Code, response.Body.String())
	}
}

func TestLogoutClearsCookie(t *testing.T) {
	sessions := &sessionApplicationStub{authenticatedID: "customer-1"}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/api/customer-sessions/current", nil)
	request.AddCookie(&http.Cookie{Name: "customer_session", Value: "raw-session-token"})
	newCustomerRouter(&customerApplicationStub{}, sessions).ServeHTTP(response, request)

	if response.Code != http.StatusNoContent || sessions.logoutRawToken != "raw-session-token" {
		t.Fatalf("unexpected logout response: %d %s", response.Code, response.Body.String())
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 || cookies[0].MaxAge >= 0 {
		t.Fatalf("expected an expired session cookie, got %+v", cookies)
	}
}
