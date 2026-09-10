package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gin-api/internal/domain"

	"github.com/go-sql-driver/mysql"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) CreateWithCardAndSession(
	ctx context.Context,
	name string,
	phone string,
	tokenHash string,
	expiresAt time.Time,
) (domain.Customer, domain.StampCard, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Customer{}, domain.StampCard{}, fmt.Errorf("begin customer registration: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO customers (name, phone)
		VALUES (?, ?)
	`, name, phone); err != nil {
		if isDuplicateEntry(err) {
			return domain.Customer{}, domain.StampCard{}, domain.ErrPhoneAlreadyExists
		}
		return domain.Customer{}, domain.StampCard{}, fmt.Errorf("insert customer: %w", err)
	}

	customer, err := findCustomerByPhone(ctx, tx, phone)
	if err != nil {
		return domain.Customer{}, domain.StampCard{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO stamp_cards (customer_id)
		VALUES (?)
	`, customer.ID); err != nil {
		return domain.Customer{}, domain.StampCard{}, fmt.Errorf("insert stamp card: %w", err)
	}
	card, err := getCardByCustomerID(ctx, tx, customer.ID)
	if err != nil {
		return domain.Customer{}, domain.StampCard{}, err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO customer_sessions (customer_id, token_hash, expires_at)
		VALUES (?, ?, ?)
	`, customer.ID, tokenHash, expiresAt.UTC()); err != nil {
		return domain.Customer{}, domain.StampCard{}, fmt.Errorf("insert registration session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return domain.Customer{}, domain.StampCard{}, fmt.Errorf("commit customer registration: %w", err)
	}
	return customer, card, nil
}

func (r *CustomerRepository) FindByPhone(ctx context.Context, phone string) (domain.Customer, error) {
	return findCustomerByPhone(ctx, r.db, phone)
}

func (r *CustomerRepository) FindByID(ctx context.Context, customerID string) (domain.Customer, error) {
	var customer domain.Customer
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, phone, created_at, updated_at
		FROM customers
		WHERE id = ?
	`, customerID).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Phone,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Customer{}, domain.ErrCustomerNotFound
	}
	if err != nil {
		return domain.Customer{}, fmt.Errorf("find customer by id: %w", err)
	}
	return customer, nil
}

func (r *CustomerRepository) GetCard(ctx context.Context, customerID string) (domain.StampCard, error) {
	return getCardByCustomerID(ctx, r.db, customerID)
}

type rowQuerier interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func findCustomerByPhone(ctx context.Context, db rowQuerier, phone string) (domain.Customer, error) {
	var customer domain.Customer
	err := db.QueryRowContext(ctx, `
		SELECT id, name, phone, created_at, updated_at
		FROM customers
		WHERE phone = ?
	`, phone).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Phone,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Customer{}, domain.ErrCustomerNotFound
	}
	if err != nil {
		return domain.Customer{}, fmt.Errorf("find customer by phone: %w", err)
	}
	return customer, nil
}

func getCardByCustomerID(ctx context.Context, db rowQuerier, customerID string) (domain.StampCard, error) {
	var card domain.StampCard
	err := db.QueryRowContext(ctx, `
		SELECT id, customer_id, stamp_count, required_stamps, created_at, updated_at
		FROM stamp_cards
		WHERE customer_id = ?
	`, customerID).Scan(
		&card.ID,
		&card.CustomerID,
		&card.StampCount,
		&card.RequiredStamps,
		&card.CreatedAt,
		&card.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.StampCard{}, domain.ErrCustomerNotFound
	}
	if err != nil {
		return domain.StampCard{}, fmt.Errorf("find customer card: %w", err)
	}
	return card, nil
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
