package mocks

import (
	"context"

	"github.com/iamviniciuss/casino-transactions/internal/core"
	"github.com/stretchr/testify/mock"
)

type mockTransactionRepository struct {
	mock.Mock
}

func NewMockTransactionRepository() *mockTransactionRepository {
	return &mockTransactionRepository{}
}

func (m *mockTransactionRepository) Save(ctx context.Context, transaction core.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *mockTransactionRepository) FindByID(ctx context.Context, transactionID string) (core.Transaction, error) {
	args := m.Called(ctx, transactionID)
	return args.Get(0).(core.Transaction), args.Error(1)
}
