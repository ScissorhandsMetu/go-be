package models

type DatabaseAppointment struct {
	ID                  int    `json:"id"`
	BarberID            int    `json:"barber_id"`            // Matches 'barber_id' in DB
	CustomerName        string `json:"customer_name"`        // Matches 'customer_name' in DB
	CustomerEmail       string `json:"customer_email"`       // Matches 'customer_email' in DB
	AppointmentDate     string `json:"appointment_date"`     // Matches 'appointment_date' in DB (ISO string)
	SlotTime            string `json:"slot_time"`            // Matches 'slot_time' in DB (HH:mm:ss)
	Status              string `json:"status"`               // Matches 'status' in DB
	VerificationToken   string `json:"verification_token"`   // Matches 'verification_token' in DB
	VerificationExpires string `json:"verification_expires"` // Matches 'verification_expires' in DB (ISO string)
}
