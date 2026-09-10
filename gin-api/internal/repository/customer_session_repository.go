package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gin-api/internal/domain"
)

type CustomerSessionRepository struct {
	db *sql.DB
}

func NewCustomerSessionRepository(db *sql.DB) *CustomerSessionRepository {
	return &CustomerSessionRepository{db: db}
}

func (r *CustomerSessionRepository) Create(
	ctx context.Context,
	customerID string,
	tokenHash string,
	expiresAt time.Time,
) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO customer_sessions (customer_id, token_hash, expires_at)
		VALUES (?, ?, ?)
	`, customerID, tokenHash, expiresAt.UTC())
	if err != nil {
		return fmt.Errorf("create customer session: %w", err)
	}
	return nil
}

func (r *CustomerSessionRepository) FindActiveCustomerID(
	ctx context.Context,
	tokenHash string,
	now time.Time,
) (string, error) {
	var customerID string
	err := r.db.QueryRowContext(ctx, `
		SELECT customer_id
		FROM customer_sessions
		WHERE token_hash = ?
		  AND revoked_at IS NULL
		  AND expires_at > ?
	`, tokenHash, now.UTC()).Scan(&customerID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", domain.ErrInvalidSession
	}
	if err != nil {
		return "", fmt.Errorf("find active customer session: %w", err)
	}
	return customerID, nil
}

func (r *CustomerSessionRepository) Revoke(ctx context.Context, tokenHash string, now time.Time) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE customer_sessions
		SET revoked_at = ?
		WHERE token_hash = ? AND revoked_at IS NULL
	`, now.UTC(), tokenHash)
	if err != nil {
		return fmt.Errorf("revoke customer session: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read revoked customer session: %w", err)
	}
	if rows == 0 {
		return domain.ErrInvalidSession
	}
	return nil
}
