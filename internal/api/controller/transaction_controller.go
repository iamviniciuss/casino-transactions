package controller

import (
	"context"
	"strconv"

	"github.com/iamviniciuss/casino-transactions/internal/api/http"
	"github.com/iamviniciuss/casino-transactions/internal/core"
)

var (
	ErrUserIDRequired = http.IntegrationError{
		StatusCode: 400,
		Message:    "user_id is required",
	}
)

type PaginatedResponse struct {
	Items  []core.Transaction `json:"items"`
	Total  int                `json:"total"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
}

type TransactionController struct {
	transactionRepository core.TransactionRepository
}

func NewTransactionController(transactionRepository core.TransactionRepository) *TransactionController {
	return &TransactionController{
		transactionRepository: transactionRepository,
	}
}

func (ctr *TransactionController) GetTransactions(ctx context.Context, m map[string]string, b []byte, qp http.QueryParams, lf http.LocalsFunc) (interface{}, *http.IntegrationError) {
	userID := string(qp.GetParam("user_id"))
	if userID == "" {
		return nil, &ErrUserIDRequired
	}

	txType := string(qp.GetParam("transaction_type"))
	limit, _ := strconv.Atoi(string(qp.GetParam("limit")))
	offset, _ := strconv.Atoi(string(qp.GetParam("offset")))

	filter := core.TransactionFilter{
		UserID:    userID,
		Type:      txType,
		Limit:     limit,
		Offset:    offset,
	}

	items, total, err := ctr.transactionRepository.FindByFilter(context.Background(), filter)
	if err != nil {
		return nil, &http.IntegrationError{
			StatusCode: 400,
			Message:    err.Error(),
		}
	}
	return PaginatedResponse{
		Items: items,
		Total: total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
