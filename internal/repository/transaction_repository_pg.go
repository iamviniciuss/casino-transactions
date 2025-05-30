package repository

import (
	"context"
	"database/sql"
	"fmt"

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

func (r *TransactionRepository) FindByFilter(ctx context.Context, f core.TransactionFilter) ([]core.Transaction, int, error) {
	baseQuery := `FROM transactions WHERE user_id = $1`
	args := []interface{}{f.UserID}
	i := 2

	if f.Type != "" {
		baseQuery += fmt.Sprintf(" AND transaction_type = $%d", i)
		args = append(args, f.Type)
		i++
	}

	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf("SELECT id, user_id, amount, transaction_type, timestamp %s ORDER BY timestamp DESC LIMIT $%d OFFSET $%d", baseQuery, i, i+1)
	args = append(args, f.Limit, f.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	transactions := []core.Transaction{}
	for rows.Next() {
		var tx core.Transaction
		if err := rows.Scan(&tx.ID, &tx.UserID, &tx.Amount, &tx.Type, &tx.Timestamp); err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, total, nil
}
