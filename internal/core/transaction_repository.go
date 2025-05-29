package core

import "context"

type TransactionRepository interface {
	Save(ctx context.Context, transaction Transaction) error
	FindByID(ctx context.Context, transaction_id string) (Transaction, error)
}
