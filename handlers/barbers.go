package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ScissorhandsMetu/go-be/db"
	"github.com/ScissorhandsMetu/go-be/models"
)

func GetBarbers(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to fetch barbers and appointments...")

	rows, err := db.DB.Query(`
        SELECT b.id, b.name, d.name AS district, b.description, b.image_url
        FROM Barbers b
        JOIN Districts d ON b.district_id = d.id
    `)
	if err != nil {
		log.Printf("Database query error: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var barbers []models.Barber
	for rows.Next() {
		var barber models.Barber
		err := rows.Scan(&barber.ID, &barber.Name, &barber.District, &barber.Description, &barber.ImageURL)
		if err != nil {
			log.Printf("Row scan error: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Fetch appointments for this barber
		appointmentsRows, err := db.DB.Query(`
            SELECT appointment_date, slot_time
            FROM Appointments
            WHERE barber_id = $1
        `, barber.ID)
		if err != nil {
			log.Printf("Database query error for appointments: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer appointmentsRows.Close()

		var appointments []models.Appointment
		for appointmentsRows.Next() {
			var appointment models.Appointment
			err := appointmentsRows.Scan(&appointment.Date, &appointment.SlotTime)
			if err != nil {
				log.Printf("Appointment row scan error: %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			appointments = append(appointments, appointment)
		}
		barber.Appointments = appointments

		barbers = append(barbers, barber)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(barbers)
	log.Println("Barbers and appointment data sent successfully.")
}
