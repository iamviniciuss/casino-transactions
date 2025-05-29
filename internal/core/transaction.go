package core

import (
	"time"
)

type TransactionType string

const (
	TransactionTypeBet TransactionType = "bet"
	TransactionTypeWin TransactionType = "win"
)

type Transaction struct {
	ID        int64           `json:"id"`
	UserID    string          `json:"user_id"`
	Type      TransactionType `json:"transaction_type"`
	Amount    float64         `json:"amount"`
	Timestamp time.Time       `json:"timestamp"`
}
