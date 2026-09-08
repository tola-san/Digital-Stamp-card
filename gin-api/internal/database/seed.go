package database

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type StaffSeed struct {
	Name     string
	Email    string
	Password string
}

// SeedStaff creates the first staff account once. Existing emails are left
// unchanged so a deployment cannot unexpectedly reset a password.
func SeedStaff(ctx context.Context, pool *DB, seed StaffSeed) (bool, error) {
	if seed.Name == "" && seed.Email == "" && seed.Password == "" {
		return false, nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(seed.Password), bcrypt.DefaultCost)
	if err != nil {
		return false, fmt.Errorf("hash staff password: %w", err)
	}

	command, err := pool.ExecContext(ctx, `
		INSERT IGNORE INTO staff (name, email, password_hash)
		VALUES (?, ?, ?)
	`, strings.TrimSpace(seed.Name), strings.ToLower(strings.TrimSpace(seed.Email)), string(passwordHash))
	if err != nil {
		return false, fmt.Errorf("seed staff account: %w", err)
	}
	rowsAffected, err := command.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("read seeded staff result: %w", err)
	}
	return rowsAffected == 1, nil
}
