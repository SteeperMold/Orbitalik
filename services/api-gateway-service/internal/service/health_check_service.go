package service

import (
	"context"
)

type HealthCheckService struct {
}

func NewHealthCheckService() *HealthCheckService {
	return &HealthCheckService{}
}

func (s *HealthCheckService) HealthCheck(ctx context.Context) error {
	return nil
}
