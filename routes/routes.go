package routes

import (
	"github.com/ScissorhandsMetu/go-be/handlers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/barbers", handlers.GetBarbers).Methods("GET")
	router.HandleFunc("/appointments", handlers.CreateAppointment).Methods("POST")
	router.HandleFunc("/appointments/{id}/status", handlers.UpdateAppointmentStatus).Methods("PUT")
	router.HandleFunc("/verify", handlers.VerifyAppointment).Methods("GET")
	router.HandleFunc("/districts", handlers.GetDistricts).Methods("GET") // Added route for districts
	router.HandleFunc("/appointments/cancel", handlers.CancelAppointment).Methods("DELETE")
}
