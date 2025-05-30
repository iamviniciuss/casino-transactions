package core

import "context"

type TransactionFilter struct {
	UserID string
	Type   string
	Limit  int
	Offset int
}

type TransactionRepository interface {
	Save(ctx context.Context, transaction Transaction) error
	FindByID(ctx context.Context, transaction_id string) (Transaction, error)
	FindByFilter(ctx context.Context, f TransactionFilter) ([]Transaction, int, error)
}
