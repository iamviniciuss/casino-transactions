package core

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Type      TransactionType `json:"transaction_type"`
	Amount    float64         `json:"amount"`
	Timestamp time.Time       `json:"timestamp"`
}

func NewTransaction(userID string, transactionType TransactionType, amount float64) (Transaction, error) {
	if !transactionType.IsValid() {
		return Transaction{}, ErrInvalidTransactionType
	}

	if amount <= 0 {
		return Transaction{}, ErrTransactionAmountZero
	}

	return Transaction{
		ID:        string(uuid.NewString()),
		UserID:    userID,
		Type:      transactionType,
		Amount:    amount,
		Timestamp: time.Now(),
	}, nil
}
