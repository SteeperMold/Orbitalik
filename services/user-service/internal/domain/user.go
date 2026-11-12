package domain

import (
	"context"

	"github.com/SteeperMold/Orbitalik/user-service/internal/models"
	"github.com/google/uuid"
)

type CreateUserParams struct {
	Email    string
	Username string
}

type UpdateUserParams struct {
	ID       uuid.UUID
	Email    string
	Username string
}

type UserService interface {
	CreateUser(ctx context.Context, params *CreateUserParams) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, params *UpdateUserParams) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, params *CreateUserParams) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, params *UpdateUserParams) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
