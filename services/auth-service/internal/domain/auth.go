package domain

import (
	"context"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
)

type AuthService interface {
	Register(ctx context.Context, email, username, password string) (*models.User, models.AccessToken, *models.RefreshToken, error)
	Login(ctx context.Context, email, password string) (*models.User, models.AccessToken, *models.RefreshToken, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.User, models.AccessToken, *models.RefreshToken, error)
	Logout(ctx context.Context, rawToken string) error
	ValidateToken(ctx context.Context, rawToken string) (*TokenValidationResult, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}
