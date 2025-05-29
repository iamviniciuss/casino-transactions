package repository

import (
	"context"
	"database/sql"

	"github.com/iamviniciuss/casino-transactions/internal/core"
	_ "github.com/lib/pq"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) Save(ctx context.Context, transaction core.Transaction) error {
	query := "INSERT INTO transactions (id, user_id, amount, transaction_type, timestamp) VALUES ($1, $2, $3, $4, $5);"

	_, err := r.db.ExecContext(ctx, query, transaction.ID, transaction.UserID, transaction.Amount, transaction.Type, transaction.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (r *TransactionRepository) FindByID(ctx context.Context, transactionID string) (core.Transaction, error) {
	var transaction core.Transaction

	query := "SELECT id, user_id, amount, transaction_type, timestamp FROM transactions WHERE id = $1;"
	row := r.db.QueryRowContext(ctx, query, transactionID)

	err := row.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Type, &transaction.Timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return core.Transaction{}, core.ErrTransactionNotFound
		}
		return core.Transaction{}, err
	}

	return transaction, nil
}
