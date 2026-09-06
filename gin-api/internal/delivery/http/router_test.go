package http_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	httpdelivery "gin-api/internal/delivery/http"
)

type healthChecker struct{ err error }

func (h healthChecker) Ping(context.Context) error { return h.err }

func TestHealth(t *testing.T) {
	router := httpdelivery.NewRouter(httpdelivery.RouterDependencies{HealthChecker: healthChecker{}})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d: %s", http.StatusOK, response.Code, response.Body.String())
	}
}

func TestHealthReportsDatabaseFailure(t *testing.T) {
	router := httpdelivery.NewRouter(httpdelivery.RouterDependencies{HealthChecker: healthChecker{err: errors.New("offline")}})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected %d, got %d: %s", http.StatusServiceUnavailable, response.Code, response.Body.String())
	}
}

func TestCORSAllowsConfiguredFrontend(t *testing.T) {
	router := httpdelivery.NewRouter(httpdelivery.RouterDependencies{HealthChecker: healthChecker{}, FrontendURL: "http://localhost:3000"})
	request := httptest.NewRequest(http.MethodOptions, "/api/health", nil)
	request.Header.Set("Origin", "http://localhost:3000")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d", http.StatusNoContent, response.Code)
	}
	if response.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Fatal("expected configured frontend origin to be allowed")
	}
}
