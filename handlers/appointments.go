package handlers

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"github.com/ScissorhandsMetu/go-be/db"
	"github.com/ScissorhandsMetu/go-be/models"
	"github.com/gorilla/mux"
)

// generateToken creates a random token for email verification.
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
	startTime := time.Now()
	log.Printf("[START] Appointment creation initiated at %s\n", startTime.Format(time.RFC3339))

	var dbAppointment models.DatabaseAppointment

	// Parse JSON body into dbAppointment struct
	if err := json.NewDecoder(r.Body).Decode(&dbAppointment); err != nil {
		log.Printf("[ERROR] Failed to parse request body: %v\n", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("[INFO] Received Appointment Data: BarberID=%d, CustomerName=%s, CustomerEmail=%s, AppointmentDate=%s, SlotTime=%s\n",
		dbAppointment.BarberID, dbAppointment.CustomerName, dbAppointment.CustomerEmail, dbAppointment.AppointmentDate, dbAppointment.SlotTime)

	// Generate a unique token for verification
	token, err := generateToken()
	if err != nil {
		log.Printf("[ERROR] Failed to generate verification token: %v\n", err)
		http.Error(w, "Failed to generate verification token", http.StatusInternalServerError)
		return
	}

	// Set expiration time for token
	expirationTime := time.Now().Add(10 * time.Minute)
	log.Printf("[INFO] Token generated and expires at: %s\n", expirationTime.Format(time.RFC3339))

	// Insert appointment into the database
	query := `
		INSERT INTO Appointments (barber_id, customer_name, customer_email, appointment_date, slot_time, status, verification_token, verification_expires)
		VALUES ($1, $2, $3, $4, $5, 'Unverified', $6, $7)
		RETURNING id, barber_id, appointment_date, slot_time;
	`
	err = db.DB.QueryRow(
		query,
		dbAppointment.BarberID,
		dbAppointment.CustomerName,
		dbAppointment.CustomerEmail,
		dbAppointment.AppointmentDate,
		dbAppointment.SlotTime,
		token,
		expirationTime,
	).Scan(&dbAppointment.ID, &dbAppointment.BarberID, &dbAppointment.AppointmentDate, &dbAppointment.SlotTime)

	if err != nil {
		log.Printf("[ERROR] Failed to insert appointment into database: %v\n", err)
		http.Error(w, "Failed to create appointment", http.StatusInternalServerError)
		return
	}

	log.Printf("[SUCCESS] Appointment created with ID=%d for BarberID=%d at %s\n", dbAppointment.ID, dbAppointment.BarberID, dbAppointment.AppointmentDate)

	// Fetch barber information
	var barberName string
	barberQuery := `SELECT name FROM Barbers WHERE id = $1;`
	err = db.DB.QueryRow(barberQuery, dbAppointment.BarberID).Scan(&barberName)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch barber information: %v\n", err)
		http.Error(w, "Failed to fetch barber information", http.StatusInternalServerError)
		return
	}
	log.Printf("[INFO] Barber Name fetched: %s\n", barberName)

	// Send verification email
	verificationLink := fmt.Sprintf("http://localhost:3001/verify?token=%s", token)
	if err := sendVerificationEmail(dbAppointment.CustomerEmail, verificationLink); err != nil {
		log.Printf("[ERROR] Failed to send verification email: %v\n", err)
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		return
	}
	log.Printf("[INFO] Verification email sent to: %s\n", dbAppointment.CustomerEmail)

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":          "Appointment created successfully",
		"barber_name":      barberName,
		"appointment_date": dbAppointment.AppointmentDate,
		"slot_time":        dbAppointment.SlotTime,
	})

	log.Printf("[END] Appointment flow completed successfully in %s\n", time.Since(startTime))
}

// VerifyAppointment handles appointment verification.
func VerifyAppointment(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	log.Printf("[START] Appointment verification initiated at %s\n", startTime.Format(time.RFC3339))

	token := r.URL.Query().Get("token")
	if token == "" {
		log.Printf("[ERROR] Missing verification token in request\n")
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE Appointments
		SET status = 'Confirmed'
		WHERE verification_token = $1 AND status = 'Unverified' AND verification_expires > NOW()
		RETURNING id;
	`

	var appointmentID int
	err := db.DB.QueryRow(query, token).Scan(&appointmentID)
	if err != nil {
		log.Printf("[ERROR] Failed to verify appointment: %v\n", err)
		http.Error(w, "Invalid or expired token", http.StatusNotFound)
		return
	}

	log.Printf("[SUCCESS] Appointment with ID=%d successfully verified.\n", appointmentID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Appointment confirmed successfully",
		"appointment_id": appointmentID,
	})

	log.Printf("[END] Appointment verification completed in %s\n", time.Since(startTime))
}

func sendVerificationEmail(toEmail, verificationLink string) error {
	from := "thescissorhandsmetu@gmail.com"
	password := "barbershop502"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte(fmt.Sprintf(
		"Subject: Appointment Verification\n\nPlease click the link to confirm your appointment: %s",
		verificationLink,
	))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Custom TLS configuration to skip certificate verification (Temporary Fix)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Disable SSL certificate validation
		ServerName:         smtpHost,
	}

	// Establish a connection to the SMTP server
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsConfig)
	if err != nil {
		log.Printf("[ERROR] Failed to establish TLS connection: %v\n", err)
		return err
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Printf("[ERROR] Failed to create SMTP client: %v\n", err)
		return err
	}
	defer client.Close()

	// Authenticate
	if err = client.Auth(auth); err != nil {
		log.Printf("[ERROR] Failed to authenticate: %v\n", err)
		return err
	}

	// Set the sender and recipient
	if err = client.Mail(from); err != nil {
		log.Printf("[ERROR] Failed to set sender: %v\n", err)
		return err
	}
	if err = client.Rcpt(toEmail); err != nil {
		log.Printf("[ERROR] Failed to set recipient: %v\n", err)
		return err
	}

	// Write the email body
	wc, err := client.Data()
	if err != nil {
		log.Printf("[ERROR] Failed to get writer: %v\n", err)
		return err
	}
	defer wc.Close()

	_, err = wc.Write(message)
	if err != nil {
		log.Printf("[ERROR] Failed to write message: %v\n", err)
		return err
	}

	// Quit SMTP session
	if err = client.Quit(); err != nil {
		log.Printf("[ERROR] Failed to close SMTP session: %v\n", err)
		return err
	}

	log.Printf("[INFO] Verification email successfully sent to %s\n", toEmail)
	return nil
}

// UpdateAppointmentStatus updates the status of an appointment.
func UpdateAppointmentStatus(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	log.Printf("[START] UpdateAppointmentStatus initiated at %s\n", startTime.Format(time.RFC3339))

	// Extract appointment ID from the URL
	vars := mux.Vars(r)
	appointmentIDStr := vars["id"]
	appointmentID, err := strconv.Atoi(appointmentIDStr)
	if err != nil {
		log.Printf("[ERROR] Invalid appointment ID format: %v\n", err)
		http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
		return
	}
	log.Printf("[INFO] Updating status for AppointmentID=%d\n", appointmentID)

	// Parse the status from the request body
	var requestData struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Printf("[ERROR] Failed to decode request body: %v\n", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("[INFO] Requested Status Update: AppointmentID=%d, NewStatus=%s\n", appointmentID, requestData.Status)

	// Update appointment status in the database
	query := `
		UPDATE Appointments
		SET status = $1
		WHERE id = $2;
	`
	result, err := db.DB.Exec(query, requestData.Status, appointmentID)
	if err != nil {
		log.Printf("[ERROR] Database error while updating status: %v\n", err)
		http.Error(w, "Failed to update appointment status", http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve rows affected: %v\n", err)
		http.Error(w, "Failed to verify update operation", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		log.Printf("[ERROR] AppointmentID=%d not found in database\n", appointmentID)
		http.Error(w, "Appointment ID not found", http.StatusNotFound)
		return
	}

	log.Printf("[SUCCESS] AppointmentID=%d successfully updated to Status=%s\n", appointmentID, requestData.Status)

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Appointment status updated successfully",
		"appointmentID": appointmentID,
		"status":        requestData.Status,
	})

	log.Printf("[END] UpdateAppointmentStatus completed in %s\n", time.Since(startTime))
}
