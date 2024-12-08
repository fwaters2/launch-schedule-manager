package launches

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateLaunch(t *testing.T) {
	store := NewInMemoryStore()
	logger := log.New(bytes.NewBuffer([]byte{}), "", log.LstdFlags)
	h := NewHandler(store, logger)

	payload := `{"mission_name":"Test Mission","launch_time":"2021-09-01T10:00:00Z","vehicle_name":"Falcon 9","launch_site":"LC-39A"}`
	req, err := http.NewRequest("POST", "/launches", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	http.HandlerFunc(h.CreateLaunch).ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rr.Code)
	}
}

func TestGetLaunchNotFound(t *testing.T) {
	store := NewInMemoryStore()
	logger := log.New(bytes.NewBuffer([]byte{}), "", log.LstdFlags)
	h := NewHandler(store, logger)

	req, err := http.NewRequest("GET", "/launches/999", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(h.GetLaunch)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}
