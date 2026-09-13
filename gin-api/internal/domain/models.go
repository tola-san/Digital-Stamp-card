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
	ID        string    `json:"id" gorm:"type:char(36);primaryKey;default:(UUID())"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	Phone     string    `json:"phone" gorm:"type:varchar(32);not null;unique"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)"`
}

func (Customer) TableName() string { return "customers" }

type Staff struct {
	ID           string    `json:"id" gorm:"type:char(36);primaryKey;default:(UUID())"`
	Name         string    `json:"name" gorm:"type:varchar(100);not null"`
	Email        string    `json:"email" gorm:"type:varchar(254);not null;unique"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)"`
}

func (Staff) TableName() string { return "staff" }

type CustomerSession struct {
	ID         string     `json:"id" gorm:"type:char(36);primaryKey;default:(UUID())"`
	CustomerID string     `json:"customer_id" gorm:"type:char(36);not null"`
	TokenHash  string     `json:"-" gorm:"type:varchar(255);not null;unique"`
	ExpiresAt  time.Time  `json:"expires_at" gorm:"type:datetime(6);not null"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty" gorm:"type:datetime(6)"`
	CreatedAt  time.Time  `json:"created_at" gorm:"type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)"`
}

func (CustomerSession) TableName() string { return "customer_sessions" }

type StaffSession struct {
	ID        string     `json:"id" gorm:"type:char(36);primaryKey;default:(UUID())"`
	StaffID   string     `json:"staff_id" gorm:"type:char(36);not null"`
	TokenHash string     `json:"-" gorm:"type:varchar(255);not null;unique"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"type:datetime(6);not null"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" gorm:"type:datetime(6)"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)"`
}

func (StaffSession) TableName() string { return "staff_sessions" }

type StampCard struct {
	ID             string    `json:"id" gorm:"type:char(36);primaryKey;default:(UUID())"`
	CustomerID     string    `json:"customer_id" gorm:"type:char(36);not null;unique"`
	StampCount     int       `json:"stamp_count" gorm:"type:int;not null;default:0"`
	RequiredStamps int       `json:"required_stamps" gorm:"type:int;not null;default:10"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)"`
}

func (StampCard) TableName() string { return "stamp_cards" }

type Reward struct {
	ID             string    `json:"id" gorm:"type:char(36);primaryKey;default:(UUID())"`
	Name           string    `json:"name" gorm:"type:varchar(255);not null"`
	Description    string    `json:"description" gorm:"type:text;not null"`
	RequiredStamps int       `json:"required_stamps" gorm:"type:int;not null"`
	Active         bool      `json:"active" gorm:"type:boolean;not null;default:true"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)"`
}

func (Reward) TableName() string { return "rewards" }

type CustomerQRToken struct {
	ID            string     `json:"id" gorm:"type:char(36);primaryKey;default:(UUID())"`
	TokenHash     string     `json:"-" gorm:"type:varchar(255);not null;unique"`
	CustomerID    string     `json:"customer_id" gorm:"type:char(36);not null"`
	Status        QRStatus   `json:"status" gorm:"type:varchar(16);not null;default:ACTIVE"`
	ExpiresAt     time.Time  `json:"expires_at" gorm:"type:datetime(6);not null"`
	UsedAt        *time.Time `json:"used_at,omitempty" gorm:"type:datetime(6)"`
	UsedByStaffID *string    `json:"used_by_staff_id,omitempty" gorm:"type:char(36)"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)"`
}

func (CustomerQRToken) TableName() string { return "customer_qr_tokens" }

type StampTransaction struct {
	ID                string          `json:"id" gorm:"type:char(36);primaryKey;default:(UUID())"`
	CustomerID        string          `json:"customer_id" gorm:"type:char(36);not null"`
	StaffID           *string         `json:"staff_id,omitempty" gorm:"type:char(36)"`
	RewardID          *string         `json:"reward_id,omitempty" gorm:"type:char(36)"`
	CustomerQRTokenID *string         `json:"customer_qr_token_id,omitempty" gorm:"type:char(36);unique"`
	Type              TransactionType `json:"type" gorm:"type:varchar(32);not null"`
	StampDelta        int             `json:"stamp_delta" gorm:"type:int;not null"`
	CreatedAt         time.Time       `json:"created_at" gorm:"type:datetime(6);not null;default:CURRENT_TIMESTAMP(6)"`
}

func (StampTransaction) TableName() string { return "stamp_transactions" }

type TransactionCursor struct {
	CreatedAt time.Time
	ID        string
}
