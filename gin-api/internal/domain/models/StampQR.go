package models

import(
	"time"
)


type QRStatus string

const (
	QRStatusActive    QRStatus = "ACTIVE"
	QRStatusUsed      QRStatus = "USED"
	QRStatusExpired   QRStatus = "EXPIRED"
	QRStatusCancelled QRStatus = "CANCELLED"
)


type StampQR struct {
	ID               string     `json:"id"`
	TokenHash        string     `json:"-"`
	StaffID          string     `json:"staff_id"`
	Status           QRStatus   `json:"status"`
	ExpiresAt        time.Time  `json:"expires_at"`
	UsedAt           *time.Time `json:"used_at,omitempty"`
	UsedByCustomerID *string    `json:"used_by_customer_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}