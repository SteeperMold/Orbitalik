package service

import (
	"context"

	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/models"
)

type TLEService struct {
	repository domain.TLERepository
}

func NewTLEService(r domain.TLERepository) *TLEService {
	return &TLEService{
		repository: r,
	}
}

func (s *TLEService) GetAllTLEs(ctx context.Context) ([]*models.TLE, error) {
	return s.repository.GetAllTLEs(ctx)
}

func (s *TLEService) GetTLEByNoradID(ctx context.Context, noradID int) (*models.TLE, error) {
	return s.repository.GetTLEByNoradID(ctx, noradID)
}

func (s *TLEService) GetTLEBySatelliteName(ctx context.Context, name string) (*models.TLE, error) {
	return s.repository.GetTLEBySatelliteName(ctx, name)
}
