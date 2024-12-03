package models

type appointments struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DistrictID   int    `json:"district_id"`
	Description  string `json:"description"`
	ImageURL     string `json:"image_url"`
	Availability []bool `json:"availability"` // Placeholder, fetched dynamically
}
