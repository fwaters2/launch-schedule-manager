package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/fwaters2/launch-schedule-manager/server/pkg/launches"
	"github.com/fwaters2/launch-schedule-manager/server/pkg/seed"
)

func InitializeInMemoryDB(store launches.Store) {
	log.Println("Initializing schema...")

	// Initialize schema (if required)
	// For in-memory stores, this might be implicit.

	log.Println("Seeding data...")

	// Seed data
	for _, launch := range seed.Launches {
		store.Create(launch)
	}

	log.Println("Data seeded successfully.")
}

func main() {
	// Basic logger
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	// In-memory store
	store := launches.NewInMemoryStore()

	// Seed data
	InitializeInMemoryDB(store)

	// Handlers
	h := launches.NewHandler(store, logger)

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/launches/{id}", h.GetLaunch).Methods("GET")
	r.HandleFunc("/launches/{id}", h.UpdateLaunch).Methods("PUT")
	r.HandleFunc("/launches/{id}", h.DeleteLaunch).Methods("DELETE")
	r.HandleFunc("/launches", h.CreateLaunch).Methods("POST")
	r.HandleFunc("/launches", h.ListLaunches).Methods("GET")

	// Start server
	addr := ":8080"
	logger.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatalf("Could not start server: %v", err)
	}
}
