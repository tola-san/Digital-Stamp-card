package domain

import "time"

type TransactionType string

const (
	TransactionStampAdded     TransactionType = "STAMP_ADDED"
	TransactionStampReversed  TransactionType = "STAMP_REVERSED"
	TransactionRewardRedeemed TransactionType = "REWARD_REDEEMED"
)

type QRStatus string

const (
	QRStatusActive    QRStatus = "ACTIVE"
	QRStatusUsed      QRStatus = "USED"
	QRStatusExpired   QRStatus = "EXPIRED"
	QRStatusCancelled QRStatus = "CANCELLED"
)

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Staff struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CustomerSession struct {
	ID         string     `json:"id"`
	CustomerID string     `json:"customer_id"`
	TokenHash  string     `json:"-"`
	ExpiresAt  time.Time  `json:"expires_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type StaffSession struct {
	ID        string     `json:"id"`
	StaffID   string     `json:"staff_id"`
	TokenHash string     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type StampCard struct {
	ID             string    `json:"id"`
	CustomerID     string    `json:"customer_id"`
	StampCount     int       `json:"stamp_count"`
	RequiredStamps int       `json:"required_stamps"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Reward struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	RequiredStamps int       `json:"required_stamps"`
	Active         bool      `json:"active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

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

type StampTransaction struct {
	ID         string          `json:"id"`
	CustomerID string          `json:"customer_id"`
	StaffID    *string         `json:"staff_id,omitempty"`
	RewardID   *string         `json:"reward_id,omitempty"`
	StampQRID  *string         `json:"stamp_qr_id,omitempty"`
	Type       TransactionType `json:"type"`
	StampDelta int             `json:"stamp_delta"`
	CreatedAt  time.Time       `json:"created_at"`
}

type TransactionCursor struct {
	CreatedAt time.Time
	ID        string
}
