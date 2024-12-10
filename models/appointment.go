package models

type DatabaseAppointment struct {
	ID              int    `json:"id"`
	BarberID        int    `json:"barber_id"`
	CustomerName    string `json:"customer_name"`
	CustomerEmail   string `json:"customer_email"`
	AppointmentDate string `json:"appointment_date"` // ISO string format
	SlotTime        string `json:"slot_time"`        // HH:mm:ss format
	Status          string `json:"status"`           // 'Pending', 'Accepted', 'Rejected'
}
