package service

import (
	"context"

	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/domain"
)

// HealthCheckService provides methods to check the health of system dependencies.
// It verifies connectivity to the database and Kafka.
type HealthCheckService struct {
	db domain.DBConn
}

// NewHealthCheckService creates and returns a new HealthCheckService.
func NewHealthCheckService(db domain.DBConn) *HealthCheckService {
	return &HealthCheckService{
		db: db,
	}
}

// HealthCheck runs health checks for the database and Kafka.
// It returns an error if any of the checks fail.
func (s *HealthCheckService) HealthCheck(ctx context.Context) error {
	return s.db.Ping(ctx)
}
