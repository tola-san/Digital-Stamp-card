package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"gin-api/internal/domain"
)

type customerStoreStub struct {
	createdName      string
	createdPhone     string
	createdTokenHash string
	createdExpiry    time.Time
	createErr        error
	customer         domain.Customer
	card             domain.StampCard
}

func (s *customerStoreStub) CreateWithCardAndSession(
	_ context.Context,
	name string,
	phone string,
	tokenHash string,
	expiresAt time.Time,
) (domain.Customer, domain.StampCard, error) {
	s.createdName = name
	s.createdPhone = phone
	s.createdTokenHash = tokenHash
	s.createdExpiry = expiresAt
	return s.customer, s.card, s.createErr
}

func (s *customerStoreStub) FindByPhone(context.Context, string) (domain.Customer, error) {
	if s.customer.ID == "" {
		return domain.Customer{}, domain.ErrCustomerNotFound
	}
	return s.customer, nil
}

func (s *customerStoreStub) FindByID(context.Context, string) (domain.Customer, error) {
	return s.customer, nil
}

func (s *customerStoreStub) GetCard(context.Context, string) (domain.StampCard, error) {
	return s.card, nil
}

type transactionStoreStub struct {
	transactions []domain.StampTransaction
	limit        int
	cursor       *domain.TransactionCursor
}

func (s *transactionStoreStub) ListByCustomerID(
	_ context.Context,
	_ string,
	limit int,
	cursor *domain.TransactionCursor,
) ([]domain.StampTransaction, error) {
	s.limit = limit
	s.cursor = cursor
	return s.transactions, nil
}

type sessionStoreStub struct {
	createdCustomerID string
	createdTokenHash  string
	activeTokenHash   string
	activeCustomerID  string
	revokedTokenHash  string
}

func (s *sessionStoreStub) Create(_ context.Context, customerID, tokenHash string, _ time.Time) error {
	s.createdCustomerID = customerID
	s.createdTokenHash = tokenHash
	return nil
}

func (s *sessionStoreStub) FindActiveCustomerID(_ context.Context, tokenHash string, _ time.Time) (string, error) {
	s.activeTokenHash = tokenHash
	if s.activeCustomerID == "" {
		return "", domain.ErrInvalidSession
	}
	return s.activeCustomerID, nil
}

func (s *sessionStoreStub) Revoke(_ context.Context, tokenHash string, _ time.Time) error {
	s.revokedTokenHash = tokenHash
	return nil
}

func TestRegisterNormalizesInputAndCreatesSession(t *testing.T) {
	fixedNow := time.Date(2026, time.September, 8, 10, 0, 0, 0, time.UTC)
	customers := &customerStoreStub{
		customer: domain.Customer{ID: "customer-1", Name: "Tola San", Phone: "+85512345678"},
		card:     domain.StampCard{ID: "card-1", CustomerID: "customer-1", RequiredStamps: 10},
	}
	service := NewCustomerService(customers, &transactionStoreStub{})
	service.now = func() time.Time { return fixedNow }

	result, err := service.Register(context.Background(), RegisterCustomerInput{
		Name: "  Tola San  ", Phone: "+855 12-345-678",
	})
	if err != nil {
		t.Fatalf("Register() returned an error: %v", err)
	}
	if customers.createdName != "Tola San" || customers.createdPhone != "+85512345678" {
		t.Fatalf("input was not normalized: name=%q phone=%q", customers.createdName, customers.createdPhone)
	}
	if result.Token == "" || len(customers.createdTokenHash) != 64 {
		t.Fatal("expected a raw session token and SHA-256 token hash")
	}
	if customers.createdTokenHash != hashSessionToken(result.Token) {
		t.Fatal("stored token hash does not match returned raw token")
	}
	if !customers.createdExpiry.Equal(fixedNow.Add(CustomerSessionDuration)) {
		t.Fatalf("unexpected expiry: %v", customers.createdExpiry)
	}
}

func TestRegisterRejectsInvalidInput(t *testing.T) {
	service := NewCustomerService(&customerStoreStub{}, &transactionStoreStub{})
	tests := []struct {
		name  string
		input RegisterCustomerInput
		want  error
	}{
		{name: "short name", input: RegisterCustomerInput{Name: "T", Phone: "+85512345678"}, want: domain.ErrInvalidName},
		{name: "local phone", input: RegisterCustomerInput{Name: "Tola", Phone: "012345678"}, want: domain.ErrInvalidPhone},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := service.Register(context.Background(), test.input)
			if !errors.Is(err, test.want) {
				t.Fatalf("expected %v, got %v", test.want, err)
			}
		})
	}
}

func TestLoginCreatesHashedSessionAndAuthenticateUsesHash(t *testing.T) {
	customers := &customerStoreStub{customer: domain.Customer{ID: "customer-1", Phone: "+85512345678"}}
	sessions := &sessionStoreStub{activeCustomerID: "customer-1"}
	service := NewCustomerSessionService(customers, sessions)

	result, err := service.Login(context.Background(), "+855 12 345 678")
	if err != nil {
		t.Fatalf("Login() returned an error: %v", err)
	}
	if sessions.createdCustomerID != "customer-1" || sessions.createdTokenHash != hashSessionToken(result.Token) {
		t.Fatal("login did not store the expected customer and token hash")
	}

	customerID, err := service.Authenticate(context.Background(), result.Token)
	if err != nil || customerID != "customer-1" {
		t.Fatalf("Authenticate() = %q, %v", customerID, err)
	}
	if sessions.activeTokenHash != sessions.createdTokenHash {
		t.Fatal("authentication did not hash the raw cookie token")
	}
}

func TestTransactionPaginationCreatesCursor(t *testing.T) {
	createdAt := time.Date(2026, time.September, 8, 10, 0, 0, 0, time.UTC)
	transactions := &transactionStoreStub{transactions: []domain.StampTransaction{
		{ID: "three", CreatedAt: createdAt},
		{ID: "two", CreatedAt: createdAt.Add(-time.Minute)},
		{ID: "one", CreatedAt: createdAt.Add(-2 * time.Minute)},
	}}
	service := NewCustomerService(&customerStoreStub{}, transactions)

	page, err := service.GetTransactions(context.Background(), "customer-1", 2, "")
	if err != nil {
		t.Fatalf("GetTransactions() returned an error: %v", err)
	}
	if transactions.limit != 3 || len(page.Transactions) != 2 || page.NextCursor == "" {
		t.Fatalf("unexpected page: %+v", page)
	}
	cursor, err := decodeTransactionCursor(page.NextCursor)
	if err != nil || cursor.ID != "two" {
		t.Fatalf("unexpected cursor: %+v, %v", cursor, err)
	}
}
