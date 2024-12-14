package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ScissorhandsMetu/go-be/db"
	"github.com/ScissorhandsMetu/go-be/models"
)

func GetDistricts(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to fetch districts...")

	// Query to fetch all districts
	rows, err := db.DB.Query(`
        SELECT id, name 
        FROM Districts
        ORDER BY name
    `)
	if err != nil {
		log.Printf("Database query error: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var districts []models.District
	for rows.Next() {
		var district models.District
		err := rows.Scan(&district.ID, &district.Name)
		if err != nil {
			log.Printf("Row scan error: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		districts = append(districts, district)
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(districts)
	log.Println("District data sent successfully.")
}
