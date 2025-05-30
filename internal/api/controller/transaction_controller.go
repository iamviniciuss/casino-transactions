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

const (
	defaultLimit  = 20
	defaultOffset = 0
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

	limit, offset := ctr.parsePaginationParams(qp)

	filter := core.TransactionFilter{
		UserID: userID,
		Type:   txType,
		Limit:  limit,
		Offset: offset,
	}

	items, total, err := ctr.transactionRepository.FindByFilter(context.Background(), filter)
	if err != nil {
		return nil, &http.IntegrationError{
			StatusCode: 400,
			Message:    err.Error(),
		}
	}
	return PaginatedResponse{
		Items:  items,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (ctr *TransactionController) parsePaginationParams(qp http.QueryParams) (limit int, offset int) {
	limitStr := string(qp.GetParam("limit"))
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	offsetStr := string(qp.GetParam("offset"))
	offset, err = strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = defaultOffset
	}

	return limit, offset
}
