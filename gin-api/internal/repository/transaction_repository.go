package repository

import (
	"context"
	"database/sql"
	"fmt"

	"gin-api/internal/domain"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) ListByCustomerID(
	ctx context.Context,
	customerID string,
	limit int,
	cursor *domain.TransactionCursor,
) ([]domain.StampTransaction, error) {
	query := `
		SELECT id, customer_id, staff_id, reward_id, stamp_qr_id, type, stamp_delta, created_at
		FROM stamp_transactions
		WHERE customer_id = ?`
	args := []any{customerID}
	if cursor != nil {
		query += ` AND (created_at < ? OR (created_at = ? AND id < ?))`
		args = append(args, cursor.CreatedAt.UTC(), cursor.CreatedAt.UTC(), cursor.ID)
	}
	query += ` ORDER BY created_at DESC, id DESC LIMIT ?`
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list customer transactions: %w", err)
	}
	defer rows.Close()

	transactions := make([]domain.StampTransaction, 0, limit)
	for rows.Next() {
		var transaction domain.StampTransaction
		var staffID, rewardID, stampQRID sql.NullString
		if err := rows.Scan(
			&transaction.ID,
			&transaction.CustomerID,
			&staffID,
			&rewardID,
			&stampQRID,
			&transaction.Type,
			&transaction.StampDelta,
			&transaction.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan customer transaction: %w", err)
		}
		transaction.StaffID = nullableStringPointer(staffID)
		transaction.RewardID = nullableStringPointer(rewardID)
		transaction.StampQRID = nullableStringPointer(stampQRID)
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate customer transactions: %w", err)
	}
	return transactions, nil
}

func nullableStringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
