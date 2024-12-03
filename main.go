package main

import (
	"log"
	"net/http"

	"github.com/ScissorhandsMetu/go-be/db"
	"github.com/ScissorhandsMetu/go-be/routes"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Connect to database
	db.Connect()

	// Set up router
	router := mux.NewRouter()
	routes.RegisterRoutes(router)

	// Enable CORS
	handler := cors.Default().Handler(router)

	// Start server
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
