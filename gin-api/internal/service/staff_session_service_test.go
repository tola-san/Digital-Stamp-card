package service


/* import packages's golang



*/
import (
	"context"
	"errors"
	"testing"
	"time"

	"gin-api/internal/domain"

	"golang.org/x/crypto/bcrypt"
)


// 
type staffStoreStub struct {
	staff          domain.Staff
	findByEmailErr error
	requestedEmail string
	requestedID    string
}


//  function Revicer from 
func (s *staffStoreStub) FindByEmail(_ context.Context, email string) (domain.Staff, error) {
	s.requestedEmail = email
	return s.staff, s.findByEmailErr
}

func (s *staffStoreStub) FindByID(_ context.Context, staffID  string) (domain.Staff, error) {
	s.requestedID = staffID
	return s.staff, nil
}

type staffSessionStoreStub struct {
	createdStaffID   string
	createdTokenHash string
	createdExpiry    time.Time
	activeStaffID    string
	activeTokenHash  string
	revokedTokenHash string
}

func (s *staffSessionStoreStub) Create(
	_ context.Context,
	staffID string,
	tokenHash string,
	expiresAt time.Time,
) error {
	s.createdStaffID = staffID
	s.createdTokenHash = tokenHash
	s.createdExpiry = expiresAt
	return nil
}

func (s *staffSessionStoreStub) FindActiveStaffID(
	_ context.Context,
	tokenHash string,
	_ time.Time,
) (string, error) {
	s.activeTokenHash = tokenHash
	if s.activeStaffID == "" {
		return "", domain.ErrInvalidStaffSession
	}
	return s.activeStaffID, nil
}

func (s *staffSessionStoreStub) Revoke(_ context.Context, tokenHash string, _ time.Time) error {
	s.revokedTokenHash = tokenHash
	return nil
}

func TestStaffLoginVerifiesPasswordAndCreatesHashedSession(t *testing.T) {

     
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	fixedNow := time.Date(2026, time.September, 12, 10, 0, 0, 0, time.UTC)
	staff := &staffStoreStub{staff: domain.Staff{
		ID: "staff-1", Name: "Shop Owner", Email: "owner@example.com", PasswordHash: string(passwordHash),
	}}
	sessions := &staffSessionStoreStub{}
	service := NewStaffSessionService(staff, sessions)
	service.now = func() time.Time { return fixedNow }

	result, err := service.Login(context.Background(), StaffLoginInput{
		Email: " OWNER@EXAMPLE.COM ", Password: "correct-password",
	})
	if err != nil {
		t.Fatalf("Login() returned an error: %v", err)
	}
	if staff.requestedEmail != "owner@example.com" {
		t.Fatalf("expected normalized email, got %q", staff.requestedEmail)
	}
	if result.Staff.ID != "staff-1" || result.Token == "" {
		t.Fatalf("unexpected login result: %+v", result)
	}
	if sessions.createdStaffID != "staff-1" || sessions.createdTokenHash == result.Token {
		t.Fatalf("session was not stored securely: %+v", sessions)
	}
	if sessions.createdTokenHash != hashSessionToken(result.Token) {
		t.Fatal("stored staff token hash does not match the returned token")
	}
	if !sessions.createdExpiry.Equal(fixedNow.Add(StaffSessionDuration)) {
		t.Fatalf("unexpected session expiry: %s", sessions.createdExpiry)
	}
}

func TestStaffLoginRejectsInvalidCredentials(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	tests := []struct {
		name     string
		staff    domain.Staff
		findErr  error
		password string
	}{
		{name: "unknown email", findErr: domain.ErrStaffNotFound, password: "any-password"},
		{name: "wrong password", staff: domain.Staff{ID: "staff-1", PasswordHash: string(passwordHash)}, password: "wrong-password"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sessions := &staffSessionStoreStub{}
			service := NewStaffSessionService(
				&staffStoreStub{staff: test.staff, findByEmailErr: test.findErr},
				sessions,
			)

			_, err := service.Login(context.Background(), StaffLoginInput{
				Email: "owner@example.com", Password: test.password,
			})
			if !errors.Is(err, domain.ErrInvalidCredentials) {
				t.Fatalf("expected invalid credentials, got %v", err)
			}
			if sessions.createdStaffID != "" {
				t.Fatal("invalid login created a staff session")
			}
		})
	}
}

func TestStaffAuthenticateAndLogoutHashRawToken(t *testing.T) {
	sessions := &staffSessionStoreStub{activeStaffID: "staff-1"}
	service := NewStaffSessionService(&staffStoreStub{}, sessions)

	staffID, err := service.Authenticate(context.Background(), "raw-token")
	if err != nil || staffID != "staff-1" {
		t.Fatalf("Authenticate() = %q, %v", staffID, err)
	}
	if sessions.activeTokenHash != hashSessionToken("raw-token") {
		t.Fatal("Authenticate() did not hash the raw token")
	}

	if err := service.Logout(context.Background(), "raw-token"); err != nil {
		t.Fatalf("Logout() returned an error: %v", err)
	}
	if sessions.revokedTokenHash != hashSessionToken("raw-token") {
		t.Fatal("Logout() did not hash the raw token")
	}
}
