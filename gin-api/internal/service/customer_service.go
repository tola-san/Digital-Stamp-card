package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"gin-api/internal/domain"
)

const CustomerSessionDuration = 30 * 24 * time.Hour

var internationalPhonePattern = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)

type CustomerStore interface {
	CreateWithCardAndSession(context.Context, string, string, string, time.Time) (domain.Customer, domain.StampCard, error)
	FindByPhone(context.Context, string) (domain.Customer, error)
	FindByID(context.Context, string) (domain.Customer, error)
	GetCard(context.Context, string) (domain.StampCard, error)
}

type TransactionStore interface {
	ListByCustomerID(context.Context, string, int, *domain.TransactionCursor) ([]domain.StampTransaction, error)
}

type RegisterCustomerInput struct {
	Name  string
	Phone string
}

type RegistrationResult struct {
	Customer  domain.Customer
	Card      domain.StampCard
	Token     string
	ExpiresAt time.Time
}

type TransactionPage struct {
	Transactions []domain.StampTransaction
	NextCursor   string
}

type CustomerService struct {
	customers    CustomerStore
	transactions TransactionStore
	now          func() time.Time
}

func NewCustomerService(customers CustomerStore, transactions TransactionStore) *CustomerService {
	return &CustomerService{customers: customers, transactions: transactions, now: time.Now}
}

func (s *CustomerService) Register(ctx context.Context, input RegisterCustomerInput) (RegistrationResult, error) {
	name := strings.TrimSpace(input.Name)
	phone := normalizePhone(input.Phone)
	if count := utf8.RuneCountInString(name); count < 2 || count > 100 {
		return RegistrationResult{}, domain.ErrInvalidName
	}
	if !internationalPhonePattern.MatchString(phone) {
		return RegistrationResult{}, domain.ErrInvalidPhone
	}

	rawToken, tokenHash, err := newSessionToken()
	if err != nil {
		return RegistrationResult{}, fmt.Errorf("generate registration session: %w", err)
	}
	expiresAt := s.now().UTC().Add(CustomerSessionDuration)
	customer, card, err := s.customers.CreateWithCardAndSession(ctx, name, phone, tokenHash, expiresAt)
	if err != nil {
		return RegistrationResult{}, err
	}
	return RegistrationResult{Customer: customer, Card: card, Token: rawToken, ExpiresAt: expiresAt}, nil
}

func (s *CustomerService) GetCustomer(ctx context.Context, customerID string) (domain.Customer, error) {
	return s.customers.FindByID(ctx, customerID)
}

func (s *CustomerService) GetCard(ctx context.Context, customerID string) (domain.StampCard, error) {
	return s.customers.GetCard(ctx, customerID)
}

func (s *CustomerService) GetTransactions(
	ctx context.Context,
	customerID string,
	limit int,
	cursor string,
) (TransactionPage, error) {
	if limit == 0 {
		limit = 20
	}
	if limit < 1 || limit > 100 {
		return TransactionPage{}, domain.ErrInvalidPageLimit
	}
	parsedCursor, err := decodeTransactionCursor(cursor)
	if err != nil {
		return TransactionPage{}, err
	}
	transactions, err := s.transactions.ListByCustomerID(ctx, customerID, limit+1, parsedCursor)
	if err != nil {
		return TransactionPage{}, err
	}
	if transactions == nil {
		transactions = make([]domain.StampTransaction, 0)
	}
	page := TransactionPage{Transactions: transactions}
	if len(transactions) > limit {
		page.Transactions = transactions[:limit]
		last := page.Transactions[len(page.Transactions)-1]
		page.NextCursor = encodeTransactionCursor(last.CreatedAt, last.ID)
	}
	return page, nil
}

func normalizePhone(phone string) string {
	replacer := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "", ".", "")
	return replacer.Replace(strings.TrimSpace(phone))
}

func encodeTransactionCursor(createdAt time.Time, id string) string {
	value := createdAt.UTC().Format(time.RFC3339Nano) + "|" + id
	return base64.RawURLEncoding.EncodeToString([]byte(value))
}

func decodeTransactionCursor(value string) (*domain.TransactionCursor, error) {
	if value == "" {
		return nil, nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil, domain.ErrInvalidCursor
	}
	parts := strings.SplitN(string(decoded), "|", 2)
	if len(parts) != 2 || parts[1] == "" {
		return nil, domain.ErrInvalidCursor
	}
	createdAt, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return nil, domain.ErrInvalidCursor
	}
	return &domain.TransactionCursor{CreatedAt: createdAt, ID: parts[1]}, nil
}
