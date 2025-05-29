package router

import (
	"context"

	"github.com/iamviniciuss/casino-transactions/internal/api/http"
)

func DataSourceRouter(httpService http.HttpService) {
	httpService.Get("/health", func(ctx context.Context, m map[string]string, b []byte, qp http.QueryParams, lf http.LocalsFunc) (interface{}, *http.IntegrationError) {
		return "OK", nil
	})
}
