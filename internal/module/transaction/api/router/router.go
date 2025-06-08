package router

import (
	controller "github.com/iamviniciuss/casino-transactions/internal/module/transaction/api/http"
	"github.com/iamviniciuss/casino-transactions/internal/module/transaction/core"
	"github.com/iamviniciuss/casino-transactions/pkg/shared/http"
)

func DataSourceRouter(httpService http.HttpService, transactionRepository core.TransactionRepository) {
	httpService.Get("/health", controller.NewHealthCheckController().Check)
	httpService.Get("/transactions", controller.NewTransactionController(transactionRepository).GetTransactions)
}
