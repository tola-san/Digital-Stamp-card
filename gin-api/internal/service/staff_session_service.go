package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gin-api/internal/domain"

	"golang.org/x/crypto/bcrypt"
)


// Define sesesion login if over 12hour expration 
const StaffSessionDuration = 1 * time.Hour


//   Dummy password any staff account
var dummyStaffPasswordHash, _ = bcrypt.GenerateFromPassword(
	[]byte("invalid-staff-password"),
	bcrypt.DefaultCost,
)

type StaffStore interface {
	FindByEmail(context.Context, string) (domain.Staff, error)
	FindByID(context.Context, string) (domain.Staff, error)
}

type StaffSessionStore interface {
	Create(context.Context, string, string, time.Time) error
	FindActiveStaffID(context.Context, string, time.Time) (string, error)
	Revoke(context.Context, string, time.Time) error
}

type StaffLoginInput struct {
	Email    string
	Password string
}

type StaffLoginResult struct {
	Staff     domain.Staff
	Token     string
	ExpiresAt time.Time
}

type StaffSessionService struct {
	staff    StaffStore
	sessions StaffSessionStore
	now      func() time.Time
}

func NewStaffSessionService(staff StaffStore, sessions StaffSessionStore) *StaffSessionService {
	return &StaffSessionService{staff: staff, sessions: sessions, now: time.Now}
}

func (s *StaffSessionService) Login(ctx context.Context, input StaffLoginInput) (StaffLoginResult, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	staff, err := s.staff.FindByEmail(ctx, email)
	passwordHash := dummyStaffPasswordHash

	if err == nil {
		passwordHash = []byte(staff.PasswordHash)
	} else if !errors.Is(err, domain.ErrStaffNotFound) {
		return StaffLoginResult{}, err
	}

	passwordErr := bcrypt.CompareHashAndPassword(passwordHash, []byte(input.Password))
	if errors.Is(err, domain.ErrStaffNotFound) || passwordErr != nil {
		return StaffLoginResult{}, domain.ErrInvalidCredentials
	}

	rawToken, tokenHash, err := newSessionToken()
	if err != nil {
		return StaffLoginResult{}, fmt.Errorf("generate staff session: %w", err)
	}
	expiresAt := s.now().UTC().Add(StaffSessionDuration)
	if err := s.sessions.Create(ctx, staff.ID, tokenHash, expiresAt); err != nil {
		return StaffLoginResult{}, err
	}

	return StaffLoginResult{Staff: staff, Token: rawToken, ExpiresAt: expiresAt}, nil
}

func (s *StaffSessionService) Authenticate(ctx context.Context, rawToken string) (string, error) {
	if rawToken == "" {
		return "", domain.ErrInvalidStaffSession
	}
	return s.sessions.FindActiveStaffID(ctx, hashSessionToken(rawToken), s.now().UTC())
}

func (s *StaffSessionService) Logout(ctx context.Context, rawToken string) error {
	if rawToken == "" {
		return domain.ErrInvalidStaffSession
	}
	return s.sessions.Revoke(ctx, hashSessionToken(rawToken), s.now().UTC())
}

func (s *StaffSessionService) GetStaff(ctx context.Context, staffID string) (domain.Staff, error) {
	return s.staff.FindByID(ctx, staffID)
}
