package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/fwaters2/launch-schedule-manager/server/pkg/launches"
)

func main() {
	// Basic logger
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	// In-memory store
	store := launches.NewInMemoryStore()

	// Handlers
	h := launches.NewHandler(store, logger)

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/launches", h.CreateLaunch).Methods("POST")
	r.HandleFunc("/launches", h.ListLaunches).Methods("GET")
	r.HandleFunc("/launches/{id}", h.GetLaunch).Methods("GET")
	r.HandleFunc("/launches/{id}", h.UpdateLaunch).Methods("PUT")
	r.HandleFunc("/launches/{id}", h.DeleteLaunch).Methods("DELETE")

	// Start server
	addr := ":8080"
	logger.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatalf("Could not start server: %v", err)
	}
}
