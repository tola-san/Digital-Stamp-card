package models

import "time"


type StampCard struct {
	ID             string    `json:"id"`
	CustomerID     string    `json:"customer_id"`
	StampCount     int       `json:"stamp_count"`
	RequiredStamps int       `json:"required_stamps"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
