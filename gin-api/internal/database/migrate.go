package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"gin-api/internal/domain"
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
		reverseQRAlreadyRenamed := false
		if name == "000005_reverse_qr_flow.up.sql" {
			reverseQRAlreadyRenamed, err = prepareReverseQRFlow(ctx, conn)
			if err != nil {
				return fmt.Errorf("prepare migration %s: %w", name, err)
			}
		}
		for _, statement := range strings.Split(string(script), ";") {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}
			if reverseQRAlreadyRenamed && strings.HasPrefix(strings.ToUpper(statement), "RENAME TABLE STAMP_QRS") {
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

	// Versioned SQL owns changes to existing tables because MySQL and MariaDB
	// report some foreign-key indexes differently to GORM. AutoMigrate is kept
	// as a safe fallback only for a genuinely missing model table; it must not
	// rewrite a healthy versioned schema during application startup.
	models := []any{
		&domain.Customer{},
		&domain.Staff{},
		&domain.CustomerSession{},
		&domain.StaffSession{},
		&domain.StampCard{},
		&domain.Reward{},
		&domain.CustomerQRToken{},
		&domain.StampTransaction{},
	}
	orm := pool.ORM.WithContext(ctx)
	for _, model := range models {
		if orm.Migrator().HasTable(model) {
			continue
		}
		if err := orm.AutoMigrate(model); err != nil {
			return fmt.Errorf("auto-migrate missing GORM model: %w", err)
		}
	}
	return nil
}

// prepareReverseQRFlow removes related foreign keys by inspecting metadata
// instead of assuming MySQL/MariaDB generated the same constraint names. It
// also recognizes a previous run that committed the table rename before a
// later DDL statement failed.
func prepareReverseQRFlow(ctx context.Context, conn *sql.Conn) (bool, error) {
	rows, err := conn.QueryContext(ctx, `
		SELECT table_name, constraint_name
		FROM information_schema.key_column_usage
		WHERE constraint_schema = DATABASE()
		  AND referenced_table_name IS NOT NULL
		  AND (
			(table_name = 'stamp_transactions' AND column_name IN ('stamp_qr_id', 'customer_qr_token_id'))
			OR
			(table_name IN ('stamp_qrs', 'customer_qr_tokens') AND column_name IN ('staff_id', 'customer_id', 'used_by_customer_id', 'used_by_staff_id'))
		  )
	`)
	if err != nil {
		return false, fmt.Errorf("inspect QR foreign keys: %w", err)
	}
	type foreignKey struct{ table, name string }
	var keys []foreignKey
	for rows.Next() {
		var key foreignKey
		if err := rows.Scan(&key.table, &key.name); err != nil {
			rows.Close()
			return false, fmt.Errorf("scan QR foreign key: %w", err)
		}
		keys = append(keys, key)
	}
	if err := rows.Close(); err != nil {
		return false, fmt.Errorf("close QR foreign key query: %w", err)
	}
	for _, key := range keys {
		query := fmt.Sprintf("ALTER TABLE `%s` DROP FOREIGN KEY `%s`", strings.ReplaceAll(key.table, "`", "``"), strings.ReplaceAll(key.name, "`", "``"))
		if _, err := conn.ExecContext(ctx, query); err != nil {
			return false, fmt.Errorf("drop foreign key %s.%s: %w", key.table, key.name, err)
		}
	}

	var oldTable, newTable int
	if err := conn.QueryRowContext(ctx, `
		SELECT
			EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'stamp_qrs'),
			EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'customer_qr_tokens')
	`).Scan(&oldTable, &newTable); err != nil {
		return false, fmt.Errorf("inspect QR tables: %w", err)
	}
	if oldTable == 0 && newTable == 0 {
		return false, fmt.Errorf("neither stamp_qrs nor customer_qr_tokens exists")
	}
	if oldTable == 1 && newTable == 1 {
		return false, fmt.Errorf("both stamp_qrs and customer_qr_tokens exist; manual reconciliation is required")
	}
	return newTable == 1, nil
}
