package repository

import (
	"context"
	"testing"
	"time"

	"github.com/iamviniciuss/casino-transactions/internal/core"
	"github.com/iamviniciuss/casino-transactions/pkg/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestPostgresRepository_Save(t *testing.T) {
	dbConn, teardown := test_utils.SetupPostgres(t)
	defer teardown()

	repo := NewTransactionRepository(dbConn)

	t.Run("save transaction", func(t *testing.T) {
		tx := core.Transaction{
			ID:        "d9e3d73b-084e-4ea7-bb36-3b382cb9f8f2",
			UserID:    "d9e3d73b-084e-4ea7-bb36-3b382cb9f8f1",
			Type:      core.TransactionTypeBet,
			Amount:    100.0,
			Timestamp: time.Now(),
		}

		err := repo.Save(context.Background(), tx)
		assert.NoError(t, err)

		tx2, err := repo.FindByID(context.Background(), tx.ID)

		assert.NoError(t, err)
		assert.Equal(t, tx.ID, tx2.ID)
		assert.Equal(t, tx.UserID, tx2.UserID)
		assert.Equal(t, tx.Type, tx2.Type)
		assert.Equal(t, tx.Amount, tx2.Amount)
		assert.Equal(t, tx.Timestamp.UTC(), tx2.Timestamp.UTC())
	})

	t.Run("save transaction with invalid amount", func(t *testing.T) {
		tx := core.Transaction{
			ID:        "d9e3d73b-084e-4ea7-bb36-3b382cb9f8g2",
			UserID:    "d9e3d73b-084e-4ea7-bb36-3b382cb9f8g1",
			Type:      core.TransactionTypeBet,
			Amount:    -100.0,
			Timestamp: time.Now(),
		}

		err := repo.Save(context.Background(), tx)
		assert.Error(t, err)
	})

	t.Run("find a invalid transaction", func(t *testing.T) {
		_, err := repo.FindByID(context.Background(), "319d51b7-1636-43f6-a4ad-10bbc388580c")
		assert.Error(t, err)
		assert.Equal(t, core.ErrTransactionNotFound, err)
	})

}
