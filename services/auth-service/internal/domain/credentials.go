package domain

import (
	"context"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
)

type CredentialsRepository interface {
	Create(ctx context.Context, cred *models.Credentials) error
	GetByUserID(ctx context.Context, userID int) (*models.Credentials, error)
	UpdatePasswordHash(ctx context.Context, userID int, hash string) error
}
