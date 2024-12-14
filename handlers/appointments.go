package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"regexp"
	"strconv"
	"time"

	"github.com/ScissorhandsMetu/go-be/db"
	"github.com/ScissorhandsMetu/go-be/models"
	"github.com/gorilla/mux"
)

// to validate email format
func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// to validate e-mail address
func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

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

	// Validate email format
	if !isValidEmail(dbAppointment.CustomerEmail) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Generate a unique token for the appointment
	token, err := generateToken()
	if err != nil {
		http.Error(w, "Failed to generate verification token", http.StatusInternalServerError)
		return
	}

	// Calculate expiration time (current time + 10 minutes)
	expirationTime := time.Now().Add(10 * time.Minute)

	// Insert appointment into database
	query := `
		INSERT INTO Appointments (barber_id, customer_name, customer_email, appointment_date, slot_time, status, token, verification_expires_at)
		VALUES ($1, $2, $3, $4, $5, 'Unverified', $6, $7)
		RETURNING id;
	`
	err = db.DB.QueryRow(query, dbAppointment.BarberID, dbAppointment.CustomerName, dbAppointment.CustomerEmail, dbAppointment.AppointmentDate, dbAppointment.SlotTime, token, expirationTime).Scan(&dbAppointment.ID)
	if err != nil {
		log.Printf("Error inserting appointment: %v\n", err)
		http.Error(w, "Failed to create appointment", http.StatusInternalServerError)
		return
	}

	fmt.Println("sending e-mail")
	// Send verification email
	verificationLink := fmt.Sprintf("http://localhost:3001/verify?token=%s", token)
	sendVerificationEmail(dbAppointment.CustomerEmail, verificationLink)

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Verification email sent",
		"appointment": dbAppointment,
	})
}

func VerifyAppointment(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	// Query to check token and expiration time
	query := `
		UPDATE Appointments
		SET status = 'Confirmed'
		WHERE token = $1 AND status = 'Unverified' AND verification_expires_at > NOW()
		RETURNING id;
	`

	var appointmentID int
	err := db.DB.QueryRow(query, token).Scan(&appointmentID)
	if err != nil {
		log.Printf("Error verifying appointment: %v\n", err)
		http.Error(w, "Invalid or expired token", http.StatusNotFound)
		return
	}

	// Respond with success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Appointment confirmed successfully",
	})
}
func sendVerificationEmail(toEmail, verificationLink string) {
	from := "thescissorhandsmetu@gmail.com"
	password := "barbershop502"
	smtpHost := "smtp.example.com"
	smtpPort := "587"

	message := []byte(fmt.Sprintf(
		"Subject: Appointment Verification\n\nPlease click the link to confirm your appointment: %s",
		verificationLink,
	))

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
	if err != nil {
		log.Printf("Error sending email: %v\n", err)
	}
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
