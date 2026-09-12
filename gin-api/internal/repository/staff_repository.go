package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gin-api/internal/domain"
)

type StaffRepository struct {
	db *sql.DB
}

func NewStaffRepository(db *sql.DB) *StaffRepository {
	return &StaffRepository{db: db}
}

func (r *StaffRepository) FindByEmail(ctx context.Context, email string) (domain.Staff, error) {
	var staff domain.Staff
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM staff
		WHERE LOWER(email) = ?
		LIMIT 1
	`, email).Scan(
		&staff.ID,
		&staff.Name,
		&staff.Email,
		&staff.PasswordHash,
		&staff.CreatedAt,
		&staff.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Staff{}, domain.ErrStaffNotFound
	}
	if err != nil {
		return domain.Staff{}, fmt.Errorf("find staff by email: %w", err)
	}
	return staff, nil
}

func (r *StaffRepository) FindByID(ctx context.Context, staffID string) (domain.Staff, error) {
	var staff domain.Staff
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM staff
		WHERE id = ?
	`, staffID).Scan(
		&staff.ID,
		&staff.Name,
		&staff.Email,
		&staff.PasswordHash,
		&staff.CreatedAt,
		&staff.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Staff{}, domain.ErrStaffNotFound
	}
	if err != nil {
		return domain.Staff{}, fmt.Errorf("find staff by id: %w", err)
	}
	return staff, nil
}
