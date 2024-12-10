package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ScissorhandsMetu/go-be/db"
	"github.com/ScissorhandsMetu/go-be/models"
	"github.com/gorilla/mux"
)

// CreateAppointment handles new appointment creation.
func CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var dbAppointment models.DatabaseAppointment

	// Parse JSON body into dbAppointment struct
	if err := json.NewDecoder(r.Body).Decode(&dbAppointment); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if dbAppointment.BarberID == 0 ||
		dbAppointment.CustomerName == "" ||
		dbAppointment.CustomerEmail == "" ||
		dbAppointment.AppointmentDate == "" ||
		dbAppointment.SlotTime == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Insert appointment into database
	query := `
		INSERT INTO Appointments (barber_id, customer_name, customer_email, appointment_date, slot_time, status)
		VALUES ($1, $2, $3, $4, $5, 'Pending')
		RETURNING id;
	`
	err := db.DB.QueryRow(query, dbAppointment.BarberID, dbAppointment.CustomerName, dbAppointment.CustomerEmail, dbAppointment.AppointmentDate, dbAppointment.SlotTime).Scan(&dbAppointment.ID)
	if err != nil {
		log.Printf("Error inserting appointment: %v\n", err)
		http.Error(w, "Failed to create appointment", http.StatusInternalServerError)
		return
	}

	// Mock response instead of email notification
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Appointment created successfully",
		"appointment": dbAppointment,
	})
}

func UpdateAppointmentStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appointmentIDStr := vars["id"]
	appointmentID, err := strconv.Atoi(appointmentIDStr)
	if err != nil {
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}

	var requestData struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Update status in the database
	query := `
		UPDATE Appointments
		SET status = $1
		WHERE id = $2;
	`
	result, err := db.DB.Exec(query, requestData.Status, appointmentID)
	if err != nil {
		log.Printf("Error updating appointment status: %v\n", err)
		http.Error(w, "Failed to update appointment", http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Appointment ID not found", http.StatusNotFound)
		return
	}

	// Mock response instead of email notification
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Appointment status updated successfully",
	})
}
