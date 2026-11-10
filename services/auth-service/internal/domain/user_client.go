package domain

import (
	"context"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
)

type UserClient interface {
	CreateUser(ctx context.Context, email, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
}
