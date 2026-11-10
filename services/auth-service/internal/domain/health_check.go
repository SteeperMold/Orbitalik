package domain

import "context"

type HealthCheckService interface {
	HealthCheck(ctx context.Context) error
}
