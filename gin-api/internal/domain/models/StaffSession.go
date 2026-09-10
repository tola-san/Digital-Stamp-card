package models


import "time"


type StaffSession struct {
	ID        string     `json:"id"`
	StaffID   string     `json:"staff_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}