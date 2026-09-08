package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"gin-api/migrations"
)

// ApplyMigrations applies each embedded migration exactly once. MySQL DDL can
// commit implicitly, so a connection-level lock serializes service startups.
func ApplyMigrations(ctx context.Context, pool *DB) error {
	conn, err := pool.Conn(ctx)
	if err != nil {
		return fmt.Errorf("reserve migration connection: %w", err)
	}
	defer conn.Close()

	var locked sql.NullInt64
	if err := conn.QueryRowContext(ctx, "SELECT GET_LOCK('digital_stamp_migrations', 10)").Scan(&locked); err != nil {
		return fmt.Errorf("lock migrations: %w", err)
	}
	if !locked.Valid || locked.Int64 != 1 {
		return fmt.Errorf("lock migrations: timed out")
	}
	defer conn.ExecContext(context.Background(), "SELECT RELEASE_LOCK('digital_stamp_migrations')")

	if _, err := conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id VARCHAR(255) PRIMARY KEY,
			applied_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
		)
	`); err != nil {
		return fmt.Errorf("create migration ledger: %w", err)
	}

	for _, name := range migrations.Up {
		var applied int
		if err := conn.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE id = ?)", name).Scan(&applied); err != nil {
			return fmt.Errorf("check migration %s: %w", name, err)
		}
		if applied == 1 {
			continue
		}

		script, err := migrations.ReadUp(name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		for _, statement := range strings.Split(string(script), ";") {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}
			if _, err := conn.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("apply migration %s: %w", name, err)
			}
		}
		if _, err := conn.ExecContext(ctx, "INSERT INTO schema_migrations (id) VALUES (?)", name); err != nil {
			return fmt.Errorf("record migration %s: %w", name, err)
		}
	}
	return nil
}
