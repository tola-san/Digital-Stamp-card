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
		if name == "000005_reverse_qr_flow.up.sql" {
			if err := applyReverseQRFlow(ctx, conn); err != nil {
				return fmt.Errorf("apply migration %s: %w", name, err)
			}
			if _, err := conn.ExecContext(ctx, "INSERT INTO schema_migrations (id) VALUES (?)", name); err != nil {
				return fmt.Errorf("record migration %s: %w", name, err)
			}
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

func applyReverseQRFlow(ctx context.Context, conn *sql.Conn) error {
	alreadyRenamed, err := prepareReverseQRFlow(ctx, conn)
	if err != nil {
		return err
	}
	if !alreadyRenamed {
		if _, err := conn.ExecContext(ctx, "RENAME TABLE stamp_qrs TO customer_qr_tokens"); err != nil {
			return fmt.Errorf("rename stamp_qrs: %w", err)
		}
	}

	if exists, err := columnExists(ctx, conn, "customer_qr_tokens", "staff_id"); err != nil {
		return err
	} else if exists {
		if _, err := conn.ExecContext(ctx, "ALTER TABLE customer_qr_tokens CHANGE COLUMN staff_id customer_id CHAR(36) CHARACTER SET ascii NOT NULL"); err != nil {
			return fmt.Errorf("rename QR customer column: %w", err)
		}
	}
	if exists, err := columnExists(ctx, conn, "customer_qr_tokens", "used_by_customer_id"); err != nil {
		return err
	} else if exists {
		if _, err := conn.ExecContext(ctx, "ALTER TABLE customer_qr_tokens CHANGE COLUMN used_by_customer_id used_by_staff_id CHAR(36) CHARACTER SET ascii"); err != nil {
			return fmt.Errorf("rename QR staff usage column: %w", err)
		}
	}

	for _, index := range []string{"stamp_qrs_staff_created_idx", "stamp_qrs_status_expiry_idx"} {
		if exists, err := indexExists(ctx, conn, "customer_qr_tokens", index); err != nil {
			return err
		} else if exists {
			if _, err := conn.ExecContext(ctx, fmt.Sprintf("ALTER TABLE customer_qr_tokens DROP INDEX `%s`", index)); err != nil {
				return fmt.Errorf("drop legacy QR index %s: %w", index, err)
			}
		}
	}
	if exists, err := indexExists(ctx, conn, "customer_qr_tokens", "customer_qr_tokens_customer_created_idx"); err != nil {
		return err
	} else if !exists {
		if _, err := conn.ExecContext(ctx, "ALTER TABLE customer_qr_tokens ADD INDEX customer_qr_tokens_customer_created_idx (customer_id, created_at DESC)"); err != nil {
			return fmt.Errorf("add QR customer index: %w", err)
		}
	}
	if exists, err := indexExists(ctx, conn, "customer_qr_tokens", "customer_qr_tokens_status_expiry_idx"); err != nil {
		return err
	} else if !exists {
		if _, err := conn.ExecContext(ctx, "ALTER TABLE customer_qr_tokens ADD INDEX customer_qr_tokens_status_expiry_idx (status, expires_at)"); err != nil {
			return fmt.Errorf("add QR status index: %w", err)
		}
	}

	if _, err := conn.ExecContext(ctx, `ALTER TABLE customer_qr_tokens
		ADD CONSTRAINT customer_qr_tokens_customer_fk FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
		ADD CONSTRAINT customer_qr_tokens_used_by_staff_fk FOREIGN KEY (used_by_staff_id) REFERENCES staff(id),
		ADD CONSTRAINT customer_qr_tokens_usage_matches_status_check CHECK (
			(status = 'USED' AND used_at IS NOT NULL AND used_by_staff_id IS NOT NULL)
			OR (status <> 'USED' AND used_at IS NULL AND used_by_staff_id IS NULL)
		)`); err != nil {
		return fmt.Errorf("add customer QR constraints: %w", err)
	}

	if exists, err := columnExists(ctx, conn, "stamp_transactions", "stamp_qr_id"); err != nil {
		return err
	} else if exists {
		if _, err := conn.ExecContext(ctx, "ALTER TABLE stamp_transactions CHANGE COLUMN stamp_qr_id customer_qr_token_id CHAR(36) CHARACTER SET ascii"); err != nil {
			return fmt.Errorf("rename transaction QR column: %w", err)
		}
	}
	if _, err := conn.ExecContext(ctx, `ALTER TABLE stamp_transactions
		ADD CONSTRAINT stamp_transactions_customer_qr_token_fk FOREIGN KEY (customer_qr_token_id) REFERENCES customer_qr_tokens(id),
		ADD CONSTRAINT stamp_transactions_type_relationships_check CHECK (
			(type = 'STAMP_ADDED' AND stamp_delta > 0 AND reward_id IS NULL)
			OR (type = 'STAMP_REVERSED' AND stamp_delta < 0 AND reward_id IS NULL AND customer_qr_token_id IS NULL)
			OR (type = 'REWARD_REDEEMED' AND stamp_delta < 0 AND reward_id IS NOT NULL AND customer_qr_token_id IS NULL)
		)`); err != nil {
		return fmt.Errorf("add transaction QR constraints: %w", err)
	}
	return nil
}

func columnExists(ctx context.Context, conn *sql.Conn, table, column string) (bool, error) {
	var exists int
	err := conn.QueryRowContext(ctx, `SELECT EXISTS(
		SELECT 1 FROM information_schema.columns
		WHERE table_schema = DATABASE() AND table_name = ? AND column_name = ?
	)`, table, column).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("inspect column %s.%s: %w", table, column, err)
	}
	return exists == 1, nil
}

func indexExists(ctx context.Context, conn *sql.Conn, table, index string) (bool, error) {
	var exists int
	err := conn.QueryRowContext(ctx, `SELECT EXISTS(
		SELECT 1 FROM information_schema.statistics
		WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?
	)`, table, index).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("inspect index %s.%s: %w", table, index, err)
	}
	return exists == 1, nil
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

	var serverVersion string
	if err := conn.QueryRowContext(ctx, "SELECT VERSION()").Scan(&serverVersion); err != nil {
		return false, fmt.Errorf("read database version: %w", err)
	}
	checkRows, err := conn.QueryContext(ctx, `
		SELECT table_name, constraint_name
		FROM information_schema.table_constraints
		WHERE constraint_schema = DATABASE()
		  AND constraint_type = 'CHECK'
		  AND (
			(table_name IN ('stamp_qrs', 'customer_qr_tokens') AND constraint_name IN (
				'stamp_qrs_usage_matches_status_check',
				'customer_qr_tokens_usage_matches_status_check'
			))
			OR
			(table_name = 'stamp_transactions' AND constraint_name = 'stamp_transactions_type_relationships_check')
		  )
	`)
	if err != nil {
		return false, fmt.Errorf("inspect QR check constraints: %w", err)
	}
	var checks []foreignKey
	for checkRows.Next() {
		var check foreignKey
		if err := checkRows.Scan(&check.table, &check.name); err != nil {
			checkRows.Close()
			return false, fmt.Errorf("scan QR check constraint: %w", err)
		}
		checks = append(checks, check)
	}
	if err := checkRows.Close(); err != nil {
		return false, fmt.Errorf("close QR check constraint query: %w", err)
	}
	for _, check := range checks {
		dropKeyword := "DROP CHECK"
		if strings.Contains(strings.ToLower(serverVersion), "mariadb") {
			dropKeyword = "DROP CONSTRAINT"
		}
		query := fmt.Sprintf("ALTER TABLE `%s` %s `%s`", strings.ReplaceAll(check.table, "`", "``"), dropKeyword, strings.ReplaceAll(check.name, "`", "``"))
		if _, err := conn.ExecContext(ctx, query); err != nil {
			return false, fmt.Errorf("drop check constraint %s.%s: %w", check.table, check.name, err)
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
