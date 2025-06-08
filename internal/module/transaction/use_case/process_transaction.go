package use_case

import (
	"context"
	"errors"
	"log"

	"github.com/iamviniciuss/casino-transactions/internal/module/transaction/core"
)

var (
	ErrInvalidTransactionType = errors.New("invalid transaction type")
	ErrTransactionAmountZero  = errors.New("transaction amount must be greater than zero")
)

type TransactionProcessor interface {
	Process(ctx context.Context, input ProcessTransactionInput) error
}

type ProcessTransactionInput struct {
	UserID          string
	Amount          float64
	GameID          string
	TransactionType core.TransactionType
}

type ProcessTransactionUseCase struct {
	transactionRepository core.TransactionRepository
}

func NewProcessTransaction(transactionRepository core.TransactionRepository) *ProcessTransactionUseCase {
	return &ProcessTransactionUseCase{transactionRepository}
}

func (pts *ProcessTransactionUseCase) Process(ctx context.Context, tx ProcessTransactionInput) error {
	transaction, err := core.NewTransaction(tx.UserID, tx.TransactionType, tx.Amount)
	if err != nil {
		return err
	}

	err = pts.transactionRepository.Save(ctx, transaction)
	if err != nil {
		return err
	}

	log.Printf("Processing transaction: %v", transaction.ID)
	return nil
}
