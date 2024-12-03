package models

type Appointment struct {
	Date     string `json:"date"`     // Appointment date
	SlotTime string `json:"slotTime"` // Appointment slot time
}

type Barber struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	District     string        `json:"district"`
	Description  string        `json:"description"`
	ImageURL     string        `json:"image"`
	Appointments []Appointment `json:"appointments"` // Add appointments field
}
