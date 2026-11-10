package grpc_test

import (
	"context"
	"time"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	grpcserver "github.com/SteeperMold/Orbitalik/auth-service/internal/transport/grpc"
	commonlog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/stretchr/testify/mock"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Register(
	ctx context.Context,
	email string,
	username string,
	password string,
) (*models.User, models.AccessToken, *models.RefreshToken, error) {
	args := m.Called(ctx, email, username, password)

	var user *models.User
	if args.Get(0) != nil {
		user = args.Get(0).(*models.User)
	}

	var refresh *models.RefreshToken
	if args.Get(2) != nil {
		refresh = args.Get(2).(*models.RefreshToken)
	}

	return user, args.Get(1).(models.AccessToken), refresh, args.Error(3)
}

func (m *mockAuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (*models.User, models.AccessToken, *models.RefreshToken, error) {
	args := m.Called(ctx, email, password)

	var user *models.User
	if args.Get(0) != nil {
		user = args.Get(0).(*models.User)
	}

	var refresh *models.RefreshToken
	if args.Get(2) != nil {
		refresh = args.Get(2).(*models.RefreshToken)
	}

	return user, args.Get(1).(models.AccessToken), refresh, args.Error(3)
}

func (m *mockAuthService) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*models.User, models.AccessToken, *models.RefreshToken, error) {
	args := m.Called(ctx, refreshToken)

	var user *models.User
	if args.Get(0) != nil {
		user = args.Get(0).(*models.User)
	}

	var refresh *models.RefreshToken
	if args.Get(2) != nil {
		refresh = args.Get(2).(*models.RefreshToken)
	}

	return user, args.Get(1).(models.AccessToken), refresh, args.Error(3)
}

func (m *mockAuthService) Logout(
	ctx context.Context,
	refreshToken string,
) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *mockAuthService) ValidateToken(
	ctx context.Context,
	accessToken string,
) (*domain.TokenValidationResult, error) {
	args := m.Called(ctx, accessToken)

	var result *domain.TokenValidationResult
	if args.Get(0) != nil {
		result = args.Get(0).(*domain.TokenValidationResult)
	}

	return result, args.Error(1)
}

type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Info(msg string, fields ...commonlog.Field) {
	m.Called(msg, fields)
}

func (m *mockLogger) Sync() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockLogger) Error(msg string, fields ...commonlog.Field) {
	m.Called(msg, fields)
}

func newTestServer(
	svc *mockAuthService,
	logger *mockLogger,
) *grpcserver.AuthServiceServer {
	return grpcserver.NewAuthServiceServer(svc, logger)
}

func testUser() *models.User {
	return &models.User{
		ID:       42,
		Email:    "test@example.com",
		Username: "testuser",
	}
}

func testRefreshToken() *models.RefreshToken {
	return &models.RefreshToken{
		Value:     "refresh-token-value",
		Hash:      "refresh-token-hash",
		ExpiresAt: time.Now().Add(time.Hour),
	}
}
