package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type StaffSeed struct {
	Name     string
	Email    string
	Password string
}

// SeedStaff creates the first staff account once. Existing emails are left
// unchanged so a deployment cannot unexpectedly reset a password.
func SeedStaff(ctx context.Context, pool *pgxpool.Pool, seed StaffSeed) (bool, error) {
	if seed.Name == "" && seed.Email == "" && seed.Password == "" {
		return false, nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(seed.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, fmt.Errorf("hash staff password: %w", err)
	}

	command, err := pool.Exec(ctx, `
		INSERT INTO staff (name, email, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO NOTHING
	`, strings.TrimSpace(seed.Name), strings.ToLower(strings.TrimSpace(seed.Email)), string(passwordHash))
	if err != nil {
		return false, fmt.Errorf("seed staff account: %w", err)
	}
	return command.RowsAffected() == 1, nil
}
