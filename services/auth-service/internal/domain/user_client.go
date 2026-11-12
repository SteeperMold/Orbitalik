package domain

import (
	"context"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	"github.com/google/uuid"
)

type UserClient interface {
	CreateUser(ctx context.Context, email, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}
