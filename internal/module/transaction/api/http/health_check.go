package controller

import (
	"context"

	"github.com/iamviniciuss/casino-transactions/pkg/shared/http"
)

type HealthCheckController struct {
}

func NewHealthCheckController() *HealthCheckController {
	return &HealthCheckController{}
}

func (ctr *HealthCheckController) Check(ctx context.Context, m map[string]string, b []byte, qp http.QueryParams, lf http.LocalsFunc) (interface{}, *http.IntegrationError) {
	return "OK", nil
}
