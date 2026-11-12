package domain

import (
	"context"
	"time"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *models.StoredRefreshToken) error
	GetByHash(ctx context.Context, hash string) (*models.StoredRefreshToken, error)
	Rotate(ctx context.Context, oldID uuid.UUID, newToken *models.StoredRefreshToken) error
	Revoke(ctx context.Context, id uuid.UUID) error
}

type TokenManager interface {
	GenerateAccessToken(user *models.User) (models.AccessToken, error)
	ValidateAccessToken(token string) (*AccessTokenClaims, error)
	GenerateRefreshToken() (*models.RefreshToken, error)
	HashToken(token string) string
}

type AccessTokenClaims struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
}

type TokenValidationResult struct {
	Valid     bool
	User      *models.User
	ExpiresAt time.Time
}
