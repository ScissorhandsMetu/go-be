package routes

import (
	"github.com/ScissorhandsMetu/go-be/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/barbers", handlers.GetBarbers).Methods("GET")
	// Add other routes (e.g., districts, appointments)
}
