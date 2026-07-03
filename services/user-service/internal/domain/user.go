package domain

import (
	"context"

	"github.com/SteeperMold/Orbitalik/user-service/internal/models"
)

type CreateUserParams struct {
	Email        string
	Username     string
	PasswordHash string
}

type UpdateUserParams struct {
	ID           int
	Email        string
	Username     string
	PasswordHash string
}

type UserService interface {
	CreateUser(ctx context.Context, params *CreateUserParams) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, params *UpdateUserParams) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, params *CreateUserParams) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, params *UpdateUserParams) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
}
