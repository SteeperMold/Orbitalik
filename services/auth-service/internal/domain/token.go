package domain

import (
	"context"
	"time"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *models.StoredRefreshToken) error
	GetByHash(ctx context.Context, hash string) (*models.StoredRefreshToken, error)
	Rotate(ctx context.Context, oldID int, newToken *models.StoredRefreshToken) error
	Revoke(ctx context.Context, id int) error
}

type TokenManager interface {
	GenerateAccessToken(user *models.User) (models.AccessToken, error)
	ValidateAccessToken(token string) (*AccessTokenClaims, error)
	GenerateRefreshToken() (*models.RefreshToken, error)
	HashToken(token string) string
}

type AccessTokenClaims struct {
	UserID    uint32
	ExpiresAt time.Time
}

type TokenValidationResult struct {
	Valid     bool
	User      *models.User
	ExpiresAt time.Time
}
