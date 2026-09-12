package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gin-api/internal/domain"
)

type StaffSessionRepository struct {
	db *sql.DB
}

func NewStaffSessionRepository(db *sql.DB) *StaffSessionRepository {
	return &StaffSessionRepository{db: db}
}

func (r *StaffSessionRepository) Create(
	ctx context.Context,
	staffID string,
	tokenHash string,
	expiresAt time.Time,
) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO staff_sessions (staff_id, token_hash, expires_at)
		VALUES (?, ?, ?)
	`, staffID, tokenHash, expiresAt.UTC())
	if err != nil {
		return fmt.Errorf("create staff session: %w", err)
	}
	return nil
}

func (r *StaffSessionRepository) FindActiveStaffID(
	ctx context.Context,
	tokenHash string,
	now time.Time,
) (string, error) {
	var staffID string
	err := r.db.QueryRowContext(ctx, `
		SELECT staff_id
		FROM staff_sessions
		WHERE token_hash = ?
		  AND revoked_at IS NULL
		  AND expires_at > ?
	`, tokenHash, now.UTC()).Scan(&staffID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", domain.ErrInvalidStaffSession
	}
	if err != nil {
		return "", fmt.Errorf("find active staff session: %w", err)
	}
	return staffID, nil
}

func (r *StaffSessionRepository) Revoke(ctx context.Context, tokenHash string, now time.Time) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE staff_sessions
		SET revoked_at = ?
		WHERE token_hash = ? AND revoked_at IS NULL
	`, now.UTC(), tokenHash)
	if err != nil {
		return fmt.Errorf("revoke staff session: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read revoked staff session: %w", err)
	}
	if rows == 0 {
		return domain.ErrInvalidStaffSession
	}
	return nil
}
