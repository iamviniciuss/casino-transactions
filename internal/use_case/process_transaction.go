// internal/application/usecase/process_transaction.go
package use_case

import (
	"context"
	"log"
	// "github.com/iamviniciuss/casino-transactions/internal/core"
)

type TransactionProcessor interface {
	Process(ctx context.Context, input interface{}) error
}

type ProcessTransactionUseCase struct {
}

func NewProcessTransaction() *ProcessTransactionUseCase {
	return &ProcessTransactionUseCase{}
}

func (uc *ProcessTransactionUseCase) Process(ctx context.Context, tx interface{}) error {
	log.Printf("Processing transaction: %v", tx)
	return nil
}
