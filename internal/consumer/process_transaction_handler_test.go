package consumer

import (
	"testing"

	"github.com/iamviniciuss/casino-transactions/internal/core"
	"github.com/iamviniciuss/casino-transactions/internal/use_case"
	"github.com/iamviniciuss/casino-transactions/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcessTransactionHandler(t *testing.T) {
	t.Run("the handler should decode the message and call the use case successfully", func(t *testing.T) {
		transactionRepo := mocks.NewMockTransactionRepository()

		transactionRepo.On("Save", t.Context(), mock.Anything).Return(nil)
		
		handler := NewProcessTransactionHandler(use_case.NewProcessTransaction(transactionRepo))
		err := handler.Handle(t.Context(), []byte(`{
			"user_id": "ea5b70af-391c-484f-a569-73e5f77cbc6a",
			"amount": 100.0,
			"transaction_type": "bet"
		}`))

		assert.NoError(t, err)
		transactionRepo.AssertNumberOfCalls(t, "Save", 1)
		transactionRepo.
			On("Save", mock.Anything, mock.MatchedBy(func(tx core.Transaction) bool {
				assert.Equal(t, 100.0, tx.Amount)
				assert.Equal(t, "bet", tx.Type)
				assert.Equal(t, "ea5b70af-391c-484f-a569-73e5f77cbc6a", tx.UserID)
				return true
			})).
			Return(nil)
	})

	t.Run("should return an error when message is invalid", func(t *testing.T) {
		transactionRepo := mocks.NewMockTransactionRepository()

		transactionRepo.On("Save", t.Context(), mock.Anything).Return(nil)
		
		handler := NewProcessTransactionHandler(use_case.NewProcessTransaction(transactionRepo))
		err := handler.Handle(t.Context(), []byte(`{invalid json}`))

		assert.Error(t, err)
		transactionRepo.AssertNumberOfCalls(t, "Save", 0)
	})
}
