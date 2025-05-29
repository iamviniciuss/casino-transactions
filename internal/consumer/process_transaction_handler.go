package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/iamviniciuss/casino-transactions/internal/core"
	"github.com/iamviniciuss/casino-transactions/internal/use_case"
)

type ProcessTransactionHandlerInput struct {
	UserID          string `json:"user_id"`
	Amount          float64 `json:"amount"`
	TransactionType core.TransactionType `json:"transaction_type"`
}

type ProcessTransactionHandler struct {
	use_case use_case.TransactionProcessor
}

func NewProcessTransactionHandler(uc use_case.TransactionProcessor) *ProcessTransactionHandler {
	return &ProcessTransactionHandler{use_case: uc}
}

func (h *ProcessTransactionHandler) Handle(ctx context.Context, msg []byte) error {
	var input ProcessTransactionHandlerInput
	if err := json.Unmarshal(msg, &input); err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}
	return h.use_case.Process(ctx, use_case.ProcessTransactionInput{
		UserID:          input.UserID,
		Amount:          input.Amount,
		TransactionType: input.TransactionType,
	})
}