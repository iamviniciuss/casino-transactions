package core

import "errors"

var (
	ErrTransactionNotFound    error = errors.New("transaction not found")
	ErrInvalidTransactionType       = errors.New("invalid transaction type")
	ErrTransactionAmountZero        = errors.New("transaction amount must be greater than zero")
)
