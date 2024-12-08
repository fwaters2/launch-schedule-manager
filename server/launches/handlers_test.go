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

func TestGetLaunch(t *testing.T) {
	func setupTestStore() Store {
		store := NewInMemoryStore()
		// Insert some test data
		store.Create(Launch{
			MissionName: "Test Mission A",
			LaunchTime: time.Date(2024, 5, 20, 14, 0, 0, 0, time.UTC),
			VehicleName: "Falcon 9",
			LaunchSite: "LC-39A",
			Status:     "scheduled",
		})
		return store
	}

	store := setupTestStore()
	logger := log.New(bytes.NewBuffer([]byte{}), "", log.LstdFlags)
	h := NewHandler(store, logger)

	req, err := http.NewRequest("GET", "/launches/1", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(h.GetLaunch)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}
