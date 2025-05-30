package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/iamviniciuss/casino-transactions/internal/core"
	"github.com/iamviniciuss/casino-transactions/pkg/test_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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


	t.Run("find transactions by filter", func(t *testing.T) {
		userID := "123e4567-e89b-12d3-a456-426614174000"

		now := time.Now().UTC()
		transactions := []core.Transaction{
			{
				ID:        uuid.NewString(),
				UserID:    userID,
				Type:      core.TransactionTypeBet,
				Amount:    50.0,
				Timestamp: now,
			},
			{
				ID:        uuid.NewString(),
				UserID:    userID,
				Type:      core.TransactionTypeWin,
				Amount:    150.0,
				Timestamp: now.Add(-1 * time.Hour),
			},
		}

		for _, tx := range transactions {
			err := repo.Save(context.Background(), tx)
			require.NoError(t, err)
		}

		t.Run("find all by user_id", func(t *testing.T) {
			filter := core.TransactionFilter{
				UserID: userID,
				Limit:  10,
				Offset: 0,
			}

			found, total, err := repo.FindByFilter(context.Background(), filter)
			require.NoError(t, err)
			assert.Len(t, found, 2)
			assert.Equal(t, 2, total)
		})

		t.Run("filter by type", func(t *testing.T) {
			filter := core.TransactionFilter{
				UserID: userID,
				Type:   string(core.TransactionTypeWin),
				Limit:  10,
				Offset: 0,
			}

			found, total, err := repo.FindByFilter(context.Background(), filter)
			require.NoError(t, err)
			assert.Len(t, found, 1)
			assert.Equal(t, core.TransactionTypeWin, found[0].Type)
			assert.Equal(t, 1, total)
		})

		t.Run("pagination works", func(t *testing.T) {
			filter := core.TransactionFilter{
				UserID: userID,
				Limit:  1,
				Offset: 1,
			}

			found, total, err := repo.FindByFilter(context.Background(), filter)
			require.NoError(t, err)
			assert.Len(t, found, 1)
			assert.Equal(t, 2, total)
		})
	})

}
