package routes

import (
	"github.com/ScissorhandsMetu/go-be/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/barbers", handlers.GetBarbers).Methods("GET")
	router.HandleFunc("/appointments", handlers.CreateAppointment).Methods("POST")
	router.HandleFunc("/appointments/{id}/status", handlers.UpdateAppointmentStatus).Methods("PUT")
}
