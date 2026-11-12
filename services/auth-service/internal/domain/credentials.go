package domain

import (
	"context"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	"github.com/google/uuid"
)

type CredentialsRepository interface {
	Create(ctx context.Context, cred *models.Credentials) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Credentials, error)
	UpdatePasswordHash(ctx context.Context, userID uuid.UUID, hash string) error
}
