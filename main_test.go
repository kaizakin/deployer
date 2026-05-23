package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	healthHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("healthHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := rr.Body.String()
	if body != "OK" {
		t.Errorf("healthHandler returned wrong body: got %q want %q", body, "OK")
	}
}

func TestStatusEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	statusHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("statusHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("statusHandler returned wrong Content-Type: got %q want %q", contentType, "application/json")
	}

	body := rr.Body.String()
	if body == "" {
		t.Error("statusHandler returned empty body")
	}
}
