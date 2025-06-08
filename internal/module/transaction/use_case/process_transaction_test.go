package use_case

import (
	"testing"

	"github.com/iamviniciuss/casino-transactions/internal/module/transaction/core"
	"github.com/iamviniciuss/casino-transactions/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessTransactionUseCase(t *testing.T) {
	t.Run("Process valid transaction", func(t *testing.T) {
		transactionRepo := mocks.NewMockTransactionRepository()

		transactionRepo.On("Save", t.Context(), mock.Anything).Return(nil)

		err := NewProcessTransaction(transactionRepo).Process(t.Context(), ProcessTransactionInput{
			UserID:          "ea5b70af-391c-484f-a569-73e5f77cbc6a",
			Amount:          100.0,
			GameID:          "0c38005d-25b9-49f4-9d9d-3a14b921173d",
			TransactionType: core.TransactionTypeBet,
		})

		assert.NoError(t, err)
		transactionRepo.AssertNumberOfCalls(t, "Save", 1)
	})

	t.Run("return error when the transaction amount is zero", func(t *testing.T) {
		transactionRepo := mocks.NewMockTransactionRepository()

		transactionRepo.On("Save", t.Context(), mock.Anything).Return(nil)

		err := NewProcessTransaction(transactionRepo).Process(t.Context(), ProcessTransactionInput{
			UserID:          "ea5b70af-391c-484f-a569-73e5f77cbc6a",
			Amount:          0.0,
			GameID:          "0c38005d-25b9-49f4-9d9d-3a14b921173d",
			TransactionType: core.TransactionTypeBet,
		})

		assert.Error(t, err, core.ErrTransactionAmountZero)
		transactionRepo.AssertNumberOfCalls(t, "Save", 0)
	})

	t.Run("return error when the transaction amount is less than zero", func(t *testing.T) {
		transactionRepo := mocks.NewMockTransactionRepository()

		transactionRepo.On("Save", t.Context(), mock.Anything).Return(nil)

		err := NewProcessTransaction(transactionRepo).Process(t.Context(), ProcessTransactionInput{
			UserID:          "ea5b70af-391c-484f-a569-73e5f77cbc6a",
			Amount:          -100.0,
			GameID:          "0c38005d-25b9-49f4-9d9d-3a14b921173d",
			TransactionType: core.TransactionTypeBet,
		})

		assert.Error(t, err, core.ErrTransactionAmountZero)
		transactionRepo.AssertNumberOfCalls(t, "Save", 0)
	})

	t.Run("return error when occur an error to save the transaction on repository", func(t *testing.T) {
		transactionRepo := mocks.NewMockTransactionRepository()

		transactionRepo.On("Save", t.Context(), mock.Anything).Return(core.ErrTransactionNotFound)

		err := NewProcessTransaction(transactionRepo).Process(t.Context(), ProcessTransactionInput{
			UserID:          "ea5b70af-391c-484f-a569-73e5f77cbc6a",
			Amount:          100.0,
			GameID:          "0c38005d-25b9-49f4-9d9d-3a14b921173d",
			TransactionType: core.TransactionTypeBet,
		})

		assert.Error(t, err, core.ErrTransactionNotFound)
		transactionRepo.AssertNumberOfCalls(t, "Save", 1)
	})
}
