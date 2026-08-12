package service

import (
	"context"

	"github.com/SteeperMold/Orbitalik/common/go/db"
)

type HealthCheckService struct {
	db db.Conn
}

func NewHealthCheckService(db db.Conn) *HealthCheckService {
	return &HealthCheckService{
		db: db,
	}
}

func (s *HealthCheckService) HealthCheck(ctx context.Context) error {
	return s.db.Ping(ctx)
}
