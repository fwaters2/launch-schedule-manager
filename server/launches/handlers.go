package launches

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	store  Store
	logger *log.Logger
}

func NewHandler(store Store, logger *log.Logger) *Handler {
	return &Handler{
		store:  store,
		logger: logger,
	}
}

func (h *Handler) CreateLaunch(w http.ResponseWriter, r *http.Request) {
	var req LaunchCreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := ValidateLaunchRequest(req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	launchTime, _ := time.Parse(time.RFC3339, req.LaunchTime)
	newLaunch := Launch{
		MissionName: req.MissionName,
		LaunchTime:  launchTime,
		VehicleName: req.VehicleName,
		LaunchSite:  req.LaunchSite,
		Status:      req.Status,
	}

	created, err := h.store.Create(newLaunch)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, created)
}

func (h *Handler) GetLaunch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	l, err := h.store.Get(id)
	if err != nil {
		if err == ErrNotFound {
			h.respondWithError(w, http.StatusNotFound, "Launch not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, l)
}

func (h *Handler) ListLaunches(w http.ResponseWriter, r *http.Request) {
	list, err := h.store.List()
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, list)
}

func (h *Handler) UpdateLaunch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var req LaunchCreateRequest
	if len(payload) > 0 {
		if err := json.Unmarshal(payload, &req); err != nil {
			h.respondWithError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
	}

	var update Launch
	if req.MissionName != "" {
		update.MissionName = req.MissionName
	}
	if req.LaunchTime != "" {
		t, err := time.Parse(time.RFC3339, req.LaunchTime)
		if err != nil {
			h.respondWithError(w, http.StatusBadRequest, ErrInvalidTimeFormat.Error())
			return
		}
		update.LaunchTime = t
	}
	if req.VehicleName != "" {
		update.VehicleName = req.VehicleName
	}
	if req.LaunchSite != "" {
		update.LaunchSite = req.LaunchSite
	}
	if req.Status != "" {
		update.Status = req.Status
	}

	updated, err := h.store.Update(id, update)
	if err != nil {
		if err == ErrNotFound {
			h.respondWithError(w, http.StatusNotFound, "Launch not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, updated)
}

func (h *Handler) DeleteLaunch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.store.Delete(id)
	if err != nil {
		if err == ErrNotFound {
			h.respondWithError(w, http.StatusNotFound, "Launch not found")
			return
		}
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "Launch deleted"})
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.logger.Printf("error: %s", message)
	http.Error(w, message, code)
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		h.logger.Printf("error marshaling response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
