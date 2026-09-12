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

type staffApplicationStub struct {
	loginResult           service.StaffLoginResult
	loginErr              error
	authenticatedID       string
	authenticateErr       error
	authenticatedRawToken string
	logoutRawToken        string
	staff                 domain.Staff
	requestedID           string
}

func (s *staffApplicationStub) Login(
	context.Context,
	service.StaffLoginInput,
) (service.StaffLoginResult, error) {
	return s.loginResult, s.loginErr
}

func (s *staffApplicationStub) Authenticate(_ context.Context, rawToken string) (string, error) {
	s.authenticatedRawToken = rawToken
	return s.authenticatedID, s.authenticateErr
}

func (s *staffApplicationStub) Logout(_ context.Context, rawToken string) error {
	s.logoutRawToken = rawToken
	return nil
}

func (s *staffApplicationStub) GetStaff(_ context.Context, staffID string) (domain.Staff, error) {
	s.requestedID = staffID
	return s.staff, nil
}

func newStaffRouter(staff *staffApplicationStub) http.Handler {
	return httpdelivery.NewRouter(httpdelivery.RouterDependencies{
		HealthChecker: healthChecker{},
		StaffService:  staff,
	})
}

func TestStaffLoginCreatesSecureSessionCookie(t *testing.T) {
	staff := &staffApplicationStub{loginResult: service.StaffLoginResult{
		Staff: domain.Staff{
			ID: "staff-1", Name: "Shop Owner", Email: "owner@example.com", PasswordHash: "secret-hash",
		},
		Token: "raw-staff-token",
	}}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/staff-sessions", bytes.NewBufferString(
		`{"email":"owner@example.com","password":"correct-password"}`,
	))
	request.Header.Set("Content-Type", "application/json")
	newStaffRouter(staff).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}
	if strings.Contains(response.Body.String(), "secret-hash") {
		t.Fatal("staff password hash was exposed in the response")
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != "staff_session" || cookies[0].Value != "raw-staff-token" {
		t.Fatalf("unexpected staff session cookie: %+v", cookies)
	}
	if !cookies[0].HttpOnly || cookies[0].Path != "/api" {
		t.Fatalf("staff session cookie is not properly restricted: %+v", cookies[0])
	}
}

func TestStaffLoginRejectsInvalidCredentials(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/staff-sessions", bytes.NewBufferString(
		`{"email":"owner@example.com","password":"wrong-password"}`,
	))
	request.Header.Set("Content-Type", "application/json")
	newStaffRouter(&staffApplicationStub{loginErr: domain.ErrInvalidCredentials}).ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized || !strings.Contains(response.Body.String(), "INVALID_CREDENTIALS") {
		t.Fatalf("unexpected invalid login response: %d %s", response.Code, response.Body.String())
	}
}

func TestProtectedStaffRouteUsesAuthenticatedStaff(t *testing.T) {
	staff := &staffApplicationStub{
		authenticatedID: "staff-1",
		staff:           domain.Staff{ID: "staff-1", Name: "Shop Owner", Email: "owner@example.com"},
	}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/staff/me", nil)
	request.AddCookie(&http.Cookie{Name: "staff_session", Value: "raw-staff-token"})
	newStaffRouter(staff).ServeHTTP(response, request)

	if response.Code != http.StatusOK || staff.requestedID != "staff-1" {
		t.Fatalf("unexpected staff response: %d %s", response.Code, response.Body.String())
	}
	if staff.authenticatedRawToken != "raw-staff-token" {
		t.Fatalf("middleware received %q", staff.authenticatedRawToken)
	}
}

func TestProtectedStaffRouteRejectsMissingSession(t *testing.T) {
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/staff/me", nil)
	newStaffRouter(&staffApplicationStub{}).ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized || !strings.Contains(response.Body.String(), "AUTHENTICATION_REQUIRED") {
		t.Fatalf("unexpected unauthorized response: %d %s", response.Code, response.Body.String())
	}
}

func TestStaffLogoutClearsCookie(t *testing.T) {
	staff := &staffApplicationStub{authenticatedID: "staff-1"}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/api/staff-sessions/current", nil)
	request.AddCookie(&http.Cookie{Name: "staff_session", Value: "raw-staff-token"})
	newStaffRouter(staff).ServeHTTP(response, request)

	if response.Code != http.StatusNoContent || staff.logoutRawToken != "raw-staff-token" {
		t.Fatalf("unexpected logout response: %d %s", response.Code, response.Body.String())
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 || cookies[0].MaxAge >= 0 {
		t.Fatalf("expected an expired staff session cookie, got %+v", cookies)
	}
}
