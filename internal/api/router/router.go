package router

import (
	"github.com/iamviniciuss/casino-transactions/internal/api/controller"
	"github.com/iamviniciuss/casino-transactions/internal/api/http"
	"github.com/iamviniciuss/casino-transactions/internal/core"
)

func DataSourceRouter(httpService http.HttpService, transactionRepository core.TransactionRepository) {
	httpService.Get("/health", controller.NewHealthCheckController().Check)
	httpService.Get("/transactions", controller.NewTransactionController(transactionRepository).GetTransactions)
}
