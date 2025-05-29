package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionCore(t *testing.T) {
	t.Run("create an transaction", func(t *testing.T) {
		tx, err := NewTransaction("ea5b70af-391c-484f-a569-73e5f77cbc6a", TransactionTypeBet, 100.0)
		assert.NoError(t, err)

		assert.Equal(t, "ea5b70af-391c-484f-a569-73e5f77cbc6a", tx.UserID)
		assert.Equal(t, TransactionTypeBet, tx.Type)
		assert.Equal(t, 100.0, tx.Amount)
		assert.NotEmpty(t, tx.ID)
	})

	t.Run("should return error when amount is zero", func(t *testing.T) {
		_, err := NewTransaction("ea5b70af-391c-484f-a569-73e5f77cbc6a", TransactionTypeBet, 0.0)
		assert.Error(t, err, ErrTransactionAmountZero)
	})

	t.Run("should return error when amount less than zero", func(t *testing.T) {
		_, err := NewTransaction("ea5b70af-391c-484f-a569-73e5f77cbc6a", TransactionTypeBet, -10.0)
		assert.Error(t, err, ErrTransactionAmountZero)
	})

	t.Run("should return error when transaction type is invalid", func(t *testing.T) {
		_, err := NewTransaction("ea5b70af-391c-484f-a569-73e5f77cbc6a", TransactionType("invalid"), -10.0)
		assert.Error(t, err, ErrInvalidTransactionType)
	})
}
