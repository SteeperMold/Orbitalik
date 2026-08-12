package grpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/auth-service/gen/authpb"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthServiceServer_Register(t *testing.T) {
	svc := new(mockAuthService)
	logger := new(mockLogger)
	server := newTestServer(svc, logger)

	user := testUser()
	refresh := testRefreshToken()
	access := models.AccessToken("access-token")

	svc.
		On(
			"Register",
			mock.Anything,
			"test@example.com",
			"testuser",
			"password123",
		).
		Return(user, access, refresh, nil).
		Once()

	req := &authpb.RegisterRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "password123",
	}

	resp, err := server.Register(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, uint32(42), resp.User.Id)
	assert.Equal(t, "test@example.com", resp.User.Email)
	assert.Equal(t, "testuser", resp.User.Username)

	assert.Equal(t, "access-token", resp.Tokens.AccessToken)
	assert.Equal(t, "refresh-token-value", resp.Tokens.RefreshToken)

	assert.InDelta(
		t,
		int64(time.Until(refresh.ExpiresAt).Seconds()),
		resp.Tokens.ExpiresInSeconds,
		1,
	)

	svc.AssertExpectations(t)
}

func TestAuthServiceServer_Register_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode codes.Code
	}{
		{
			name:     "user not found",
			err:      domain.ErrUserNotFound,
			wantCode: codes.NotFound,
		},
		{
			name:     "user already exists",
			err:      domain.ErrUserAlreadyExists,
			wantCode: codes.AlreadyExists,
		},
		{
			name:     "invalid credentials",
			err:      domain.ErrInvalidCredentials,
			wantCode: codes.Unauthenticated,
		},
		{
			name:     "invalid token",
			err:      domain.ErrTokenInvalid,
			wantCode: codes.Unauthenticated,
		},
		{
			name:     "token expired",
			err:      domain.ErrTokenExpired,
			wantCode: codes.Unauthenticated,
		},
		{
			name:     "invalid email",
			err:      domain.ErrInvalidEmail,
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "invalid username",
			err:      domain.ErrInvalidUsername,
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "invalid password",
			err:      domain.ErrInvalidPassword,
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "weak password",
			err:      domain.ErrWeakPassword,
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "password too long",
			err:      domain.ErrPasswordTooLong,
			wantCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockAuthService)
			logger := new(mockLogger)
			server := newTestServer(svc, logger)

			svc.
				On(
					"Register",
					mock.Anything,
					"test@example.com",
					"testuser",
					"password123",
				).
				Return(
					(*models.User)(nil),
					models.AccessToken(""),
					(*models.RefreshToken)(nil),
					tt.err,
				).
				Once()

			resp, err := server.Register(
				context.Background(),
				&authpb.RegisterRequest{
					Email:    "test@example.com",
					Username: "testuser",
					Password: "password123",
				},
			)

			assert.Nil(t, resp)
			require.Error(t, err)

			assert.Equal(t, tt.wantCode, status.Code(err))

			svc.AssertExpectations(t)
		})
	}
}

func TestAuthServiceServer_Login(t *testing.T) {
	svc := new(mockAuthService)
	logger := new(mockLogger)
	server := newTestServer(svc, logger)

	user := testUser()
	refresh := testRefreshToken()
	access := models.AccessToken("access-token")

	svc.
		On(
			"Login",
			mock.Anything,
			"test@example.com",
			"password123",
		).
		Return(user, access, refresh, nil).
		Once()

	resp, err := server.Login(
		context.Background(),
		&authpb.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		},
	)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, uint32(42), resp.User.Id)
	assert.Equal(t, "test@example.com", resp.User.Email)
	assert.Equal(t, "testuser", resp.User.Username)

	assert.Equal(t, "access-token", resp.Tokens.AccessToken)
	assert.Equal(t, "refresh-token-value", resp.Tokens.RefreshToken)

	svc.AssertExpectations(t)
}

func TestAuthServiceServer_Login_Error(t *testing.T) {
	svc := new(mockAuthService)
	logger := new(mockLogger)
	server := newTestServer(svc, logger)

	svc.
		On(
			"Login",
			mock.Anything,
			"test@example.com",
			"wrong-password",
		).
		Return(
			(*models.User)(nil),
			models.AccessToken(""),
			(*models.RefreshToken)(nil),
			domain.ErrInvalidCredentials,
		).
		Once()

	resp, err := server.Login(
		context.Background(),
		&authpb.LoginRequest{
			Email:    "test@example.com",
			Password: "wrong-password",
		},
	)

	assert.Nil(t, resp)
	require.Error(t, err)

	assert.Equal(t, codes.Unauthenticated, status.Code(err))

	svc.AssertExpectations(t)
}

func TestAuthServiceServer_RefreshToken(t *testing.T) {
	svc := new(mockAuthService)
	logger := new(mockLogger)
	server := newTestServer(svc, logger)

	user := testUser()
	refresh := testRefreshToken()
	access := models.AccessToken("new-access-token")

	svc.
		On(
			"RefreshToken",
			mock.Anything,
			"old-refresh-token",
		).
		Return(user, access, refresh, nil).
		Once()

	resp, err := server.RefreshToken(
		context.Background(),
		&authpb.RefreshTokenRequest{
			RefreshToken: "old-refresh-token",
		},
	)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, uint32(42), resp.User.Id)
	assert.Equal(t, "new-access-token", resp.Tokens.AccessToken)
	assert.Equal(t, "refresh-token-value", resp.Tokens.RefreshToken)

	svc.AssertExpectations(t)
}

func TestAuthServiceServer_RefreshToken_InvalidToken(t *testing.T) {
	svc := new(mockAuthService)
	logger := new(mockLogger)
	server := newTestServer(svc, logger)

	svc.
		On(
			"RefreshToken",
			mock.Anything,
			"invalid-refresh-token",
		).
		Return(
			(*models.User)(nil),
			models.AccessToken(""),
			(*models.RefreshToken)(nil),
			domain.ErrTokenInvalid,
		).
		Once()

	resp, err := server.RefreshToken(
		context.Background(),
		&authpb.RefreshTokenRequest{
			RefreshToken: "invalid-refresh-token",
		},
	)

	assert.Nil(t, resp)
	require.Error(t, err)

	assert.Equal(t, codes.Unauthenticated, status.Code(err))

	svc.AssertExpectations(t)
}

func TestAuthServiceServer_Logout(t *testing.T) {
	svc := new(mockAuthService)
	logger := new(mockLogger)
	server := newTestServer(svc, logger)

	svc.
		On(
			"Logout",
			mock.Anything,
			"refresh-token",
		).
		Return(nil).
		Once()

	resp, err := server.Logout(
		context.Background(),
		&authpb.LogoutRequest{
			RefreshToken: "refresh-token",
		},
	)

	require.NoError(t, err)
	require.NotNil(t, resp)

	svc.AssertExpectations(t)
}

func TestAuthServiceServer_Logout_InvalidToken(t *testing.T) {
	svc := new(mockAuthService)
	logger := new(mockLogger)
	server := newTestServer(svc, logger)

	svc.
		On(
			"Logout",
			mock.Anything,
			"invalid-token",
		).
		Return(domain.ErrTokenInvalid).
		Once()

	resp, err := server.Logout(
		context.Background(),
		&authpb.LogoutRequest{
			RefreshToken: "invalid-token",
		},
	)

	assert.Nil(t, resp)
	require.Error(t, err)

	assert.Equal(t, codes.Unauthenticated, status.Code(err))

	svc.AssertExpectations(t)
}

func TestAuthServiceServer_ValidateToken_Valid(t *testing.T) {
	svc := new(mockAuthService)
	logger := new(mockLogger)
	server := newTestServer(svc, logger)

	user := testUser()
	expiresAt := time.Now().Add(time.Hour)

	result := &domain.TokenValidationResult{
		Valid:     true,
		User:      user,
		ExpiresAt: expiresAt,
	}

	svc.
		On(
			"ValidateToken",
			mock.Anything,
			"access-token",
		).
		Return(result, nil).
		Once()

	resp, err := server.ValidateToken(
		context.Background(),
		&authpb.ValidateTokenRequest{
			AccessToken: "access-token",
		},
	)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.True(t, resp.Valid)

	require.NotNil(t, resp.User)
	assert.Equal(t, uint32(42), resp.User.Id)
	assert.Equal(t, "test@example.com", resp.User.Email)
	assert.Equal(t, "testuser", resp.User.Username)

	require.NotNil(t, resp.ExpiresAt)
	assert.WithinDuration(
		t,
		expiresAt,
		resp.ExpiresAt.AsTime(),
		time.Second,
	)

	svc.AssertExpectations(t)
}

func TestAuthServiceServer_ValidateToken_Invalid(t *testing.T) {
	svc := new(mockAuthService)
	logger := new(mockLogger)
	server := newTestServer(svc, logger)

	svc.
		On("ValidateToken", mock.Anything, "invalid-access-token").
		Return(&domain.TokenValidationResult{Valid: false}, nil).
		Once()

	resp, err := server.ValidateToken(
		context.Background(),
		&authpb.ValidateTokenRequest{
			AccessToken: "invalid-access-token",
		},
	)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.False(t, resp.Valid)

	svc.AssertExpectations(t)
}

func TestAuthServiceServer_ErrorMapping(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code codes.Code
		msg  string
	}{
		{
			name: "user not found",
			err:  domain.ErrUserNotFound,
			code: codes.NotFound,
			msg:  "user not found",
		},
		{
			name: "user already exists",
			err:  domain.ErrUserAlreadyExists,
			code: codes.AlreadyExists,
			msg:  "user already exists",
		},
		{
			name: "invalid credentials",
			err:  domain.ErrInvalidCredentials,
			code: codes.Unauthenticated,
			msg:  "invalid credentials",
		},
		{
			name: "invalid token",
			err:  domain.ErrTokenInvalid,
			code: codes.Unauthenticated,
			msg:  "invalid token",
		},
		{
			name: "token expired",
			err:  domain.ErrTokenExpired,
			code: codes.Unauthenticated,
			msg:  "token expired",
		},
		{
			name: "invalid email",
			err:  domain.ErrInvalidEmail,
			code: codes.InvalidArgument,
			msg:  "invalid email",
		},
		{
			name: "invalid username",
			err:  domain.ErrInvalidUsername,
			code: codes.InvalidArgument,
			msg:  "invalid username",
		},
		{
			name: "invalid password",
			err:  domain.ErrInvalidPassword,
			code: codes.InvalidArgument,
			msg:  "invalid password",
		},
		{
			name: "weak password",
			err:  domain.ErrWeakPassword,
			code: codes.InvalidArgument,
			msg:  "password is too weak",
		},
		{
			name: "password too long",
			err:  domain.ErrPasswordTooLong,
			code: codes.InvalidArgument,
			msg:  "password is too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := new(mockAuthService)
			logger := new(mockLogger)
			server := newTestServer(svc, logger)

			svc.
				On(
					"Register",
					mock.Anything,
					"test@example.com",
					"testuser",
					"password123",
				).
				Return(
					(*models.User)(nil),
					models.AccessToken(""),
					(*models.RefreshToken)(nil),
					tt.err,
				).
				Once()

			resp, err := server.Register(
				context.Background(),
				&authpb.RegisterRequest{
					Email:    "test@example.com",
					Username: "testuser",
					Password: "password123",
				},
			)

			assert.Nil(t, resp)
			require.Error(t, err)

			st, ok := status.FromError(err)
			require.True(t, ok)

			assert.Equal(t, tt.code, st.Code())
			assert.Equal(t, tt.msg, st.Message())

			svc.AssertExpectations(t)
		})
	}
}

func TestAuthServiceServer_InternalError(t *testing.T) {
	svc := new(mockAuthService)
	logger := new(mockLogger)
	server := newTestServer(svc, logger)

	internalErr := errors.New("database connection failed")

	svc.
		On(
			"Register",
			mock.Anything,
			"test@example.com",
			"testuser",
			"password123",
		).
		Return(
			(*models.User)(nil),
			models.AccessToken(""),
			(*models.RefreshToken)(nil),
			internalErr,
		).
		Once()

	logger.
		On("Error", "internal server error", mock.Anything).
		Return().
		Once()

	resp, err := server.Register(
		context.Background(),
		&authpb.RegisterRequest{
			Email:    "test@example.com",
			Username: "testuser",
			Password: "password123",
		},
	)

	assert.Nil(t, resp)
	require.Error(t, err)

	assert.Equal(t, codes.Internal, status.Code(err))
	assert.Equal(t, "internal server error", status.Convert(err).Message())

	svc.AssertExpectations(t)
}
