package service

import (
	"context"
	"errors"
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"strings"

	"github.com/SteeperMold/Orbitalik/user-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/user-service/internal/models"
)

type UserService struct {
	repo   domain.UserRepository
	logger applog.Logger
}

func NewUserService(repo domain.UserRepository, logger applog.Logger) domain.UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserService) CreateUser(ctx context.Context, params *domain.CreateUserParams) (*models.User, error) {
	if params.Email == "" {
		return nil, domain.ErrEmailRequired
	}

	if params.Username == "" {
		return nil, domain.ErrUsernameRequired
	}

	params.Email = strings.ToLower(strings.TrimSpace(params.Email))

	existing, err := s.repo.GetUserByEmail(ctx, params.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	if existing != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	return s.repo.CreateUser(ctx, params)
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *UserService) UpdateUser(ctx context.Context, params *domain.UpdateUserParams) (*models.User, error) {
	existing, err := s.repo.GetUserByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, domain.ErrUserNotFound
	}

	if params.Email != "" {
		params.Email = strings.TrimSpace(strings.ToLower(params.Email))

		other, err := s.repo.GetUserByEmail(ctx, params.Email)
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}

		if other != nil && other.ID != params.ID {
			return nil, domain.ErrUserAlreadyExists
		}
	}

	if params.Username != "" {
		params.Username = strings.TrimSpace(params.Username)
	}

	return s.repo.UpdateUser(ctx, params)
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if user == nil {
		return domain.ErrUserNotFound
	}

	return s.repo.DeleteUser(ctx, id)
}
