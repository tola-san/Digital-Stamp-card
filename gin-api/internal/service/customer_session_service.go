package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gin-api/internal/domain"
)

type CustomerSessionStore interface {
	Create(context.Context, string, string, time.Time) error
	FindActiveCustomerID(context.Context, string, time.Time) (string, error)
	Revoke(context.Context, string, time.Time) error
}

type LoginResult struct {
	Customer  domain.Customer
	Token     string
	ExpiresAt time.Time
}

type CustomerSessionService struct {
	customers CustomerStore
	sessions  CustomerSessionStore
	now       func() time.Time
}

func NewCustomerSessionService(customers CustomerStore, sessions CustomerSessionStore) *CustomerSessionService {
	return &CustomerSessionService{customers: customers, sessions: sessions, now: time.Now}
}

func (s *CustomerSessionService) Login(ctx context.Context, phone string) (LoginResult, error) {
	phone = normalizePhone(phone)
	if !internationalPhonePattern.MatchString(phone) {
		return LoginResult{}, domain.ErrInvalidPhone
	}
	customer, err := s.customers.FindByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, domain.ErrCustomerNotFound) {
			return LoginResult{}, domain.ErrCustomerNotFound
		}
		return LoginResult{}, err
	}
	rawToken, tokenHash, err := newSessionToken()
	if err != nil {
		return LoginResult{}, fmt.Errorf("generate customer session: %w", err)
	}
	expiresAt := s.now().UTC().Add(CustomerSessionDuration)
	if err := s.sessions.Create(ctx, customer.ID, tokenHash, expiresAt); err != nil {
		return LoginResult{}, err
	}
	return LoginResult{Customer: customer, Token: rawToken, ExpiresAt: expiresAt}, nil
}

func (s *CustomerSessionService) Authenticate(ctx context.Context, rawToken string) (string, error) {
	if rawToken == "" {
		return "", domain.ErrInvalidSession
	}
	return s.sessions.FindActiveCustomerID(ctx, hashSessionToken(rawToken), s.now().UTC())
}

func (s *CustomerSessionService) Logout(ctx context.Context, rawToken string) error {
	if rawToken == "" {
		return domain.ErrInvalidSession
	}
	return s.sessions.Revoke(ctx, hashSessionToken(rawToken), s.now().UTC())
}
