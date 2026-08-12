package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testUser = &models.User{
		ID:       42,
		Email:    "john@example.com",
		Username: "john",
	}

	testAccessToken models.AccessToken = "access-token"
)

func TestValidateRegister(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		username string
		password string
		wantErr  error
	}{
		{
			name:     "valid",
			email:    "john@example.com",
			username: "john",
			password: "password123",
		},
		{
			name:     "empty email",
			email:    "",
			username: "john",
			password: "password123",
			wantErr:  domain.ErrInvalidEmail,
		},
		{
			name:     "whitespace email",
			email:    "   ",
			username: "john",
			password: "password123",
			wantErr:  domain.ErrInvalidEmail,
		},
		{
			name:     "invalid email",
			email:    "not-an-email",
			username: "john",
			password: "password123",
			wantErr:  domain.ErrInvalidEmail,
		},
		{
			name:     "empty username",
			email:    "john@example.com",
			username: "",
			password: "password123",
			wantErr:  domain.ErrInvalidUsername,
		},
		{
			name:     "whitespace username",
			email:    "john@example.com",
			username: "   ",
			password: "password123",
			wantErr:  domain.ErrInvalidUsername,
		},
		{
			name:     "username too short",
			email:    "john@example.com",
			username: "ab",
			password: "password123",
			wantErr:  domain.ErrInvalidUsername,
		},
		{
			name:     "username exactly 3 characters",
			email:    "john@example.com",
			username: "abc",
			password: "password123",
		},
		{
			name:     "username exactly 32 characters",
			email:    "john@example.com",
			username: "abcdefghijklmnopqrstuvwxyzabcdef",
			password: "password123",
		},
		{
			name:     "username too long",
			email:    "john@example.com",
			username: "abcdefghijklmnopqrstuvwxyzabcdefg",
			password: "password123",
			wantErr:  domain.ErrInvalidUsername,
		},
		{
			name:     "unicode username counted as runes",
			email:    "john@example.com",
			username: "ąžuolas",
			password: "password123",
		},
		{
			name:     "empty password",
			email:    "john@example.com",
			username: "john",
			password: "",
			wantErr:  domain.ErrInvalidPassword,
		},
		{
			name:     "password too short",
			email:    "john@example.com",
			username: "john",
			password: "1234567",
			wantErr:  domain.ErrWeakPassword,
		},
		{
			name:     "password exactly 8 characters",
			email:    "john@example.com",
			username: "john",
			password: "12345678",
		},
		{
			name:     "password exactly 128 characters",
			email:    "john@example.com",
			username: "john",
			password: string(make([]byte, 128)),
		},
		{
			name:     "password too long",
			email:    "john@example.com",
			username: "john",
			password: string(make([]byte, 129)),
			wantErr:  domain.ErrPasswordTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRegister(
				tt.email,
				tt.username,
				tt.password,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestValidateLogin(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		wantErr  error
	}{
		{
			name:     "valid",
			email:    "john@example.com",
			password: "password",
		},
		{
			name:     "empty email",
			email:    "",
			password: "password",
			wantErr:  domain.ErrInvalidEmail,
		},
		{
			name:     "empty password",
			email:    "john@example.com",
			password: "",
			wantErr:  domain.ErrInvalidPassword,
		},
		{
			name:     "both empty",
			email:    "",
			password: "",
			wantErr:  domain.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLogin(tt.email, tt.password)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, credentials, tokens, tokenManager, hasher, users := newTestService()

		ctx := context.Background()

		refresh := &models.RefreshToken{
			Value:     "refresh-token",
			Hash:      "refresh-hash",
			ExpiresAt: time.Now().Add(time.Hour),
		}

		hasher.
			On("Hash", "password123").
			Return("password-hash", nil).
			Once()

		users.
			On("CreateUser", ctx, "john@example.com", "john").
			Return(testUser, nil).
			Once()

		credentials.
			On(
				"Create",
				ctx,
				mock.MatchedBy(func(c *models.Credentials) bool {
					return c.UserID == testUser.ID &&
						c.PasswordHash == "password-hash"
				}),
			).
			Return(nil).
			Once()

		tokenManager.
			On("GenerateAccessToken", testUser).
			Return(testAccessToken, nil).
			Once()

		tokenManager.
			On("GenerateRefreshToken").
			Return(refresh, nil).
			Once()

		tokens.
			On(
				"Save",
				ctx,
				mock.MatchedBy(func(token *models.StoredRefreshToken) bool {
					return token.UserID == testUser.ID &&
						token.Hash == refresh.Hash &&
						token.ExpiresAt.Equal(refresh.ExpiresAt)
				}),
			).
			Return(nil).
			Once()

		user, access, returnedRefresh, err :=
			svc.Register(
				ctx,
				"john@example.com",
				"john",
				"password123",
			)

		require.NoError(t, err)
		require.Equal(t, testUser, user)
		require.Equal(t, testAccessToken, access)
		require.Equal(t, refresh, returnedRefresh)

		hasher.AssertExpectations(t)
		users.AssertExpectations(t)
		credentials.AssertExpectations(t)
		tokenManager.AssertExpectations(t)
		tokens.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		svc, _, _, _, _, _ := newTestService()

		user, access, refresh, err :=
			svc.Register(
				context.Background(),
				"bad-email",
				"john",
				"password123",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("hashing fails", func(t *testing.T) {
		svc, _, _, _, hasher, users := newTestService()

		expectedErr := errors.New("hashing failed")

		hasher.
			On("Hash", "password123").
			Return("", expectedErr).
			Once()

		user, access, refresh, err :=
			svc.Register(
				context.Background(),
				"john@example.com",
				"john",
				"password123",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, expectedErr)

		users.AssertNotCalled(t, "CreateUser")
	})

	t.Run("user creation fails", func(t *testing.T) {
		svc, _, _, _, hasher, users := newTestService()

		expectedErr := errors.New("user service unavailable")

		hasher.
			On("Hash", "password123").
			Return("password-hash", nil).
			Once()

		users.
			On("CreateUser", mock.Anything, "john@example.com", "john").
			Return(nil, expectedErr).
			Once()

		user, access, refresh, err :=
			svc.Register(
				context.Background(),
				"john@example.com",
				"john",
				"password123",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("credentials creation fails", func(t *testing.T) {
		svc, credentials, _, _, hasher, users := newTestService()

		expectedErr := errors.New("credentials creation failed")

		hasher.
			On("Hash", "password123").
			Return("password-hash", nil).
			Once()

		users.
			On("CreateUser", mock.Anything, "john@example.com", "john").
			Return(testUser, nil).
			Once()

		credentials.
			On("Create", mock.Anything, mock.Anything).
			Return(expectedErr).
			Once()

		user, access, refresh, err :=
			svc.Register(
				context.Background(),
				"john@example.com",
				"john",
				"password123",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("access token generation fails", func(t *testing.T) {
		svc, credentials, tokens, tokenManager, hasher, users := newTestService()

		expectedErr := errors.New("token generation failed")

		hasher.
			On("Hash", "password123").
			Return("password-hash", nil).
			Once()

		users.
			On("CreateUser", mock.Anything, "john@example.com", "john").
			Return(testUser, nil).
			Once()

		credentials.
			On("Create", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		tokenManager.
			On("GenerateAccessToken", testUser).
			Return(models.AccessToken(""), expectedErr).
			Once()

		user, access, refresh, err := svc.Register(
			context.Background(),
			"john@example.com",
			"john",
			"password123",
		)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, expectedErr)

		tokens.AssertNotCalled(t, "Save")
	})

	t.Run("refresh token generation fails", func(t *testing.T) {
		svc, _, tokens, tokenManager, hasher, users := newTestService()

		expectedErr := errors.New("refresh token generation failed")

		hasher.
			On("Hash", "password123").
			Return("password-hash", nil).
			Once()

		users.
			On("CreateUser", mock.Anything, "john@example.com", "john").
			Return(testUser, nil).
			Once()

		svcCredentials := svc.credentials.(*MockCredentialsRepository)

		svcCredentials.
			On("Create", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		tokenManager.
			On("GenerateAccessToken", testUser).
			Return(testAccessToken, nil).
			Once()

		tokenManager.
			On("GenerateRefreshToken").
			Return(nil, expectedErr).
			Once()

		user, access, refresh, err :=
			svc.Register(
				context.Background(),
				"john@example.com",
				"john",
				"password123",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, expectedErr)

		tokens.AssertNotCalled(t, "Save")
	})

	t.Run("saving refresh token fails", func(t *testing.T) {
		svc, credentials, tokens, tokenManager, hasher, users :=
			newTestService()

		expectedErr := errors.New("database unavailable")

		refresh := &models.RefreshToken{
			Value:     "refresh-token",
			Hash:      "refresh-hash",
			ExpiresAt: time.Now().Add(time.Hour),
		}

		hasher.
			On("Hash", "password123").
			Return("password-hash", nil).
			Once()

		users.
			On("CreateUser", mock.Anything, "john@example.com", "john").
			Return(testUser, nil).
			Once()

		credentials.
			On("Create", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		tokenManager.
			On("GenerateAccessToken", testUser).
			Return(testAccessToken, nil).
			Once()

		tokenManager.
			On("GenerateRefreshToken").
			Return(refresh, nil).
			Once()

		tokens.
			On("Save", mock.Anything, mock.Anything).
			Return(expectedErr).
			Once()

		user, access, returnedRefresh, err :=
			svc.Register(
				context.Background(),
				"john@example.com",
				"john",
				"password123",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, returnedRefresh)
		require.ErrorIs(t, err, expectedErr)
	})
}

func TestAuthService_Login(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, credentials, tokens, tokenManager, hasher, users :=
			newTestService()

		ctx := context.Background()

		storedCredentials := &models.Credentials{
			UserID:       testUser.ID,
			PasswordHash: "stored-hash",
		}

		refresh := &models.RefreshToken{
			Value:     "refresh-token",
			Hash:      "refresh-hash",
			ExpiresAt: time.Now().Add(time.Hour),
		}

		users.
			On("GetByEmail", ctx, testUser.Email).
			Return(testUser, nil).
			Once()

		credentials.
			On("GetByUserID", ctx, testUser.ID).
			Return(storedCredentials, nil).
			Once()

		hasher.
			On("Compare", "stored-hash", "password123").
			Return(nil).
			Once()

		tokenManager.
			On("GenerateAccessToken", testUser).
			Return(testAccessToken, nil).
			Once()

		tokenManager.
			On("GenerateRefreshToken").
			Return(refresh, nil).
			Once()

		tokens.
			On("Save", ctx, mock.MatchedBy(func(token *models.StoredRefreshToken) bool {
				return token.UserID == testUser.ID &&
					token.Hash == refresh.Hash &&
					token.ExpiresAt.Equal(refresh.ExpiresAt)
			})).
			Return(nil).
			Once()

		user, access, returnedRefresh, err :=
			svc.Login(ctx, testUser.Email, "password123")

		require.NoError(t, err)
		require.Equal(t, testUser, user)
		require.Equal(t, testAccessToken, access)
		require.Equal(t, refresh, returnedRefresh)

		users.AssertExpectations(t)
		credentials.AssertExpectations(t)
		hasher.AssertExpectations(t)
		tokenManager.AssertExpectations(t)
		tokens.AssertExpectations(t)
	})

	t.Run("invalid email", func(t *testing.T) {
		svc, _, _, _, _, users := newTestService()

		user, access, refresh, err :=
			svc.Login(
				context.Background(),
				"",
				"password123",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, domain.ErrInvalidEmail)

		users.AssertNotCalled(t, "GetByEmail")
	})

	t.Run("empty password", func(t *testing.T) {
		svc, _, _, _, _, users := newTestService()

		user, access, refresh, err :=
			svc.Login(
				context.Background(),
				testUser.Email,
				"",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, domain.ErrInvalidPassword)

		users.AssertNotCalled(t, "GetByEmail")
	})

	t.Run("user lookup fails", func(t *testing.T) {
		svc, _, _, _, _, users := newTestService()

		expectedErr := domain.ErrUserNotFound

		users.
			On("GetByEmail", mock.Anything, testUser.Email).
			Return(nil, expectedErr).
			Once()

		user, access, refresh, err :=
			svc.Login(
				context.Background(),
				testUser.Email,
				"password123",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("credentials lookup fails", func(t *testing.T) {
		svc, credentials, _, _, _, users := newTestService()

		users.
			On("GetByEmail", mock.Anything, testUser.Email).
			Return(testUser, nil).
			Once()

		credentials.
			On("GetByUserID", mock.Anything, testUser.ID).
			Return(nil, errors.New("database error")).
			Once()

		user, access, refresh, err :=
			svc.Login(
				context.Background(),
				testUser.Email,
				"password123",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, domain.ErrInvalidCredentials)
	})

	t.Run("wrong password", func(t *testing.T) {
		svc, credentials, _, _, hasher, users := newTestService()

		users.
			On("GetByEmail", mock.Anything, testUser.Email).
			Return(testUser, nil).
			Once()

		credentials.
			On("GetByUserID", mock.Anything, testUser.ID).
			Return(&models.Credentials{
				UserID:       testUser.ID,
				PasswordHash: "stored-hash",
			}, nil).
			Once()

		hasher.
			On("Compare", "stored-hash", "wrong-password").
			Return(domain.ErrInvalidCredentials).
			Once()

		user, access, refresh, err :=
			svc.Login(
				context.Background(),
				testUser.Email,
				"wrong-password",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, domain.ErrInvalidCredentials)
	})

	t.Run("access token generation fails", func(t *testing.T) {
		svc, credentials, _, tokenManager, hasher, users := newTestService()

		expectedErr := errors.New("access token failed")

		users.
			On("GetByEmail", mock.Anything, testUser.Email).
			Return(testUser, nil).
			Once()

		credentials.
			On("GetByUserID", mock.Anything, testUser.ID).
			Return(&models.Credentials{
				UserID:       testUser.ID,
				PasswordHash: "stored-hash",
			}, nil).
			Once()

		hasher.
			On("Compare", "stored-hash", "password123").
			Return(nil).
			Once()

		tokenManager.
			On("GenerateAccessToken", testUser).
			Return(models.AccessToken(""), expectedErr).
			Once()

		user, access, refresh, err :=
			svc.Login(
				context.Background(),
				testUser.Email,
				"password123",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("refresh token generation fails", func(t *testing.T) {
		svc, credentials, _, tokenManager, hasher, users :=
			newTestService()

		expectedErr := errors.New("refresh token failed")

		users.
			On("GetByEmail", mock.Anything, testUser.Email).
			Return(testUser, nil).
			Once()

		credentials.
			On("GetByUserID", mock.Anything, testUser.ID).
			Return(&models.Credentials{
				UserID:       testUser.ID,
				PasswordHash: "stored-hash",
			}, nil).
			Once()

		hasher.
			On("Compare", "stored-hash", "password123").
			Return(nil).
			Once()

		tokenManager.
			On("GenerateAccessToken", testUser).
			Return(testAccessToken, nil).
			Once()

		tokenManager.
			On("GenerateRefreshToken").
			Return(nil, expectedErr).
			Once()

		user, access, refresh, err :=
			svc.Login(
				context.Background(),
				testUser.Email,
				"password123",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("saving refresh token fails", func(t *testing.T) {
		svc, credentials, tokens, tokenManager, hasher, users :=
			newTestService()

		expectedErr := errors.New("save failed")

		refresh := &models.RefreshToken{
			Value:     "refresh-token",
			Hash:      "refresh-hash",
			ExpiresAt: time.Now().Add(time.Hour),
		}

		users.
			On("GetByEmail", mock.Anything, testUser.Email).
			Return(testUser, nil).
			Once()

		credentials.
			On("GetByUserID", mock.Anything, testUser.ID).
			Return(&models.Credentials{
				UserID:       testUser.ID,
				PasswordHash: "stored-hash",
			}, nil).
			Once()

		hasher.
			On("Compare", "stored-hash", "password123").
			Return(nil).
			Once()

		tokenManager.
			On("GenerateAccessToken", testUser).
			Return(testAccessToken, nil).
			Once()

		tokenManager.
			On("GenerateRefreshToken").
			Return(refresh, nil).
			Once()

		tokens.
			On("Save", mock.Anything, mock.Anything).
			Return(expectedErr).
			Once()

		user, access, returnedRefresh, err :=
			svc.Login(
				context.Background(),
				testUser.Email,
				"password123",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, returnedRefresh)
		require.ErrorIs(t, err, expectedErr)
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	t.Run("empty token", func(t *testing.T) {
		svc, _, _, _, _, _ := newTestService()

		user, access, refresh, err :=
			svc.RefreshToken(context.Background(), "")

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, domain.ErrTokenInvalid)
	})

	t.Run("token not found", func(t *testing.T) {
		svc, _, tokens, _, _, _ := newTestService()

		tokens.
			On("GetByHash", mock.Anything, "token-hash").
			Return(nil, errors.New("not found")).
			Once()

		svc.tokenManager.(*MockTokenManager).
			On("HashToken", "refresh-token").
			Return("token-hash").
			Once()

		user, access, refresh, err :=
			svc.RefreshToken(
				context.Background(),
				"refresh-token",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, domain.ErrTokenInvalid)
	})

	t.Run("expired token", func(t *testing.T) {
		svc, _, tokens, tokenManager, _, _ := newTestService()

		tokenManager.
			On("HashToken", "refresh-token").
			Return("token-hash").
			Once()

		tokens.
			On("GetByHash", mock.Anything, "token-hash").
			Return(&models.StoredRefreshToken{
				ID:        1,
				UserID:    testUser.ID,
				Hash:      "token-hash",
				ExpiresAt: time.Now().Add(-time.Hour),
			}, nil).
			Once()

		user, access, refresh, err :=
			svc.RefreshToken(
				context.Background(),
				"refresh-token",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, domain.ErrTokenExpired)
	})

	t.Run("revoked token", func(t *testing.T) {
		svc, _, tokens, tokenManager, _, _ := newTestService()

		revokedAt := time.Now().Add(-time.Minute)

		tokenManager.
			On("HashToken", "refresh-token").
			Return("token-hash").
			Once()

		tokens.
			On("GetByHash", mock.Anything, "token-hash").
			Return(&models.StoredRefreshToken{
				ID:        1,
				UserID:    testUser.ID,
				Hash:      "token-hash",
				ExpiresAt: time.Now().Add(time.Hour),
				RevokedAt: &revokedAt,
			}, nil).
			Once()

		user, access, refresh, err :=
			svc.RefreshToken(
				context.Background(),
				"refresh-token",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, domain.ErrTokenInvalid)
	})

	t.Run("user lookup fails", func(t *testing.T) {
		svc, _, tokens, tokenManager, _, users :=
			newTestService()

		expectedErr := domain.ErrUserNotFound

		tokenManager.
			On("HashToken", "refresh-token").
			Return("token-hash").
			Once()

		tokens.
			On("GetByHash", mock.Anything, "token-hash").
			Return(&models.StoredRefreshToken{
				ID:        1,
				UserID:    testUser.ID,
				ExpiresAt: time.Now().Add(time.Hour),
			}, nil).
			Once()

		users.
			On("GetByID", mock.Anything, testUser.ID).
			Return(nil, expectedErr).
			Once()

		user, access, refresh, err :=
			svc.RefreshToken(
				context.Background(),
				"refresh-token",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("success", func(t *testing.T) {
		svc, _, tokens, tokenManager, _, users :=
			newTestService()

		ctx := context.Background()

		oldToken := &models.StoredRefreshToken{
			ID:        10,
			UserID:    testUser.ID,
			Hash:      "old-hash",
			ExpiresAt: time.Now().Add(time.Hour),
		}

		newRefresh := &models.RefreshToken{
			Value:     "new-refresh-token",
			Hash:      "new-refresh-hash",
			ExpiresAt: time.Now().Add(2 * time.Hour),
		}

		tokenManager.
			On("HashToken", "old-refresh-token").
			Return("old-hash").
			Once()

		tokens.
			On("GetByHash", ctx, "old-hash").
			Return(oldToken, nil).
			Once()

		users.
			On("GetByID", ctx, testUser.ID).
			Return(testUser, nil).
			Once()

		tokenManager.
			On("GenerateAccessToken", testUser).
			Return(testAccessToken, nil).
			Once()

		tokenManager.
			On("GenerateRefreshToken").
			Return(newRefresh, nil).
			Once()

		tokens.
			On(
				"Rotate",
				ctx,
				oldToken.ID,
				mock.MatchedBy(func(token *models.StoredRefreshToken) bool {
					return token.UserID == testUser.ID &&
						token.Hash == newRefresh.Hash &&
						token.ExpiresAt.Equal(newRefresh.ExpiresAt)
				}),
			).
			Return(nil).
			Once()

		user, access, refresh, err :=
			svc.RefreshToken(ctx, "old-refresh-token")

		require.NoError(t, err)
		require.Equal(t, testUser, user)
		require.Equal(t, testAccessToken, access)
		require.Equal(t, newRefresh, refresh)

		tokens.AssertExpectations(t)
		tokenManager.AssertExpectations(t)
		users.AssertExpectations(t)
	})

	t.Run("access token generation fails", func(t *testing.T) {
		svc, _, tokens, tokenManager, _, users :=
			newTestService()

		expectedErr := errors.New("access token failed")

		tokenManager.
			On("HashToken", "refresh-token").
			Return("hash").
			Once()

		tokens.
			On("GetByHash", mock.Anything, "hash").
			Return(&models.StoredRefreshToken{
				ID:        1,
				UserID:    testUser.ID,
				ExpiresAt: time.Now().Add(time.Hour),
			}, nil).
			Once()

		users.
			On("GetByID", mock.Anything, testUser.ID).
			Return(testUser, nil).
			Once()

		tokenManager.
			On("GenerateAccessToken", testUser).
			Return(models.AccessToken(""), expectedErr).
			Once()

		user, access, refresh, err :=
			svc.RefreshToken(
				context.Background(),
				"refresh-token",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, expectedErr)

		tokenManager.AssertNotCalled(t, "GenerateRefreshToken")
		tokens.AssertNotCalled(t, "Rotate")
	})

	t.Run("refresh token generation fails", func(t *testing.T) {
		svc, _, tokens, tokenManager, _, users :=
			newTestService()

		expectedErr := errors.New("refresh generation failed")

		tokenManager.
			On("HashToken", "refresh-token").
			Return("hash").
			Once()

		tokens.
			On("GetByHash", mock.Anything, "hash").
			Return(&models.StoredRefreshToken{
				ID:        1,
				UserID:    testUser.ID,
				ExpiresAt: time.Now().Add(time.Hour),
			}, nil).
			Once()

		users.
			On("GetByID", mock.Anything, testUser.ID).
			Return(testUser, nil).
			Once()

		tokenManager.
			On("GenerateAccessToken", testUser).
			Return(testAccessToken, nil).
			Once()

		tokenManager.
			On("GenerateRefreshToken").
			Return(nil, expectedErr).
			Once()

		user, access, refresh, err :=
			svc.RefreshToken(
				context.Background(),
				"refresh-token",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, expectedErr)

		tokens.AssertNotCalled(t, "Rotate")
	})

	t.Run("rotation fails", func(t *testing.T) {
		svc, _, tokens, tokenManager, _, users :=
			newTestService()

		expectedErr := errors.New("rotation failed")

		newRefresh := &models.RefreshToken{
			Value:     "new-token",
			Hash:      "new-hash",
			ExpiresAt: time.Now().Add(time.Hour),
		}

		tokenManager.
			On("HashToken", "refresh-token").
			Return("old-hash").
			Once()

		tokens.
			On("GetByHash", mock.Anything, "old-hash").
			Return(&models.StoredRefreshToken{
				ID:        1,
				UserID:    testUser.ID,
				ExpiresAt: time.Now().Add(time.Hour),
			}, nil).
			Once()

		users.
			On("GetByID", mock.Anything, testUser.ID).
			Return(testUser, nil).
			Once()

		tokenManager.
			On("GenerateAccessToken", testUser).
			Return(testAccessToken, nil).
			Once()

		tokenManager.
			On("GenerateRefreshToken").
			Return(newRefresh, nil).
			Once()

		tokens.
			On("Rotate", mock.Anything, 1, mock.Anything).
			Return(expectedErr).
			Once()

		user, access, refresh, err :=
			svc.RefreshToken(
				context.Background(),
				"refresh-token",
			)

		require.Nil(t, user)
		require.Empty(t, access)
		require.Nil(t, refresh)
		require.ErrorIs(t, err, expectedErr)
	})
}

func TestAuthService_Logout(t *testing.T) {
	t.Run("empty token", func(t *testing.T) {
		svc, _, _, _, _, _ := newTestService()

		err := svc.Logout(context.Background(), "")

		require.ErrorIs(t, err, domain.ErrTokenInvalid)
	})

	t.Run("token lookup fails", func(t *testing.T) {
		svc, _, tokens, tokenManager, _, _ := newTestService()

		tokenManager.
			On("HashToken", "refresh-token").
			Return("hash").
			Once()

		tokens.
			On("GetByHash", mock.Anything, "hash").
			Return(nil, errors.New("database error")).
			Once()

		err := svc.Logout(
			context.Background(),
			"refresh-token",
		)

		require.EqualError(t, err, "database error")

		tokens.AssertExpectations(t)
		tokenManager.AssertExpectations(t)
	})

	t.Run("already revoked", func(t *testing.T) {
		svc, _, tokens, tokenManager, _, _ := newTestService()

		revokedAt := time.Now()

		tokenManager.
			On("HashToken", "refresh-token").
			Return("hash").
			Once()

		tokens.
			On("GetByHash", mock.Anything, "hash").
			Return(&models.StoredRefreshToken{
				ID:        123,
				RevokedAt: &revokedAt,
			}, nil).
			Once()

		err := svc.Logout(
			context.Background(),
			"refresh-token",
		)

		require.NoError(t, err)

		tokens.AssertNotCalled(t, "Revoke")
	})

	t.Run("success", func(t *testing.T) {
		svc, _, tokens, tokenManager, _, _ := newTestService()

		tokenManager.
			On("HashToken", "refresh-token").
			Return("hash").
			Once()

		tokens.
			On("GetByHash", mock.Anything, "hash").
			Return(&models.StoredRefreshToken{
				ID: 123,
			}, nil).
			Once()

		tokens.
			On("Revoke", mock.Anything, 123).
			Return(nil).
			Once()

		err := svc.Logout(
			context.Background(),
			"refresh-token",
		)

		require.NoError(t, err)

		tokens.AssertExpectations(t)
	})

	t.Run("revoke fails", func(t *testing.T) {
		svc, _, tokens, tokenManager, _, _ := newTestService()

		expectedErr := errors.New("database unavailable")

		tokenManager.
			On("HashToken", "refresh-token").
			Return("hash").
			Once()

		tokens.
			On("GetByHash", mock.Anything, "hash").
			Return(&models.StoredRefreshToken{
				ID: 123,
			}, nil).
			Once()

		tokens.
			On("Revoke", mock.Anything, 123).
			Return(expectedErr).
			Once()

		err := svc.Logout(
			context.Background(),
			"refresh-token",
		)

		require.ErrorIs(t, err, expectedErr)
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	t.Run("empty token", func(t *testing.T) {
		svc, _, _, _, _, _ := newTestService()

		result, err :=
			svc.ValidateToken(
				context.Background(),
				"",
			)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.False(t, result.Valid)
	})

	t.Run("invalid token", func(t *testing.T) {
		svc, _, _, tokenManager, _, _ := newTestService()

		tokenManager.
			On("ValidateAccessToken", "invalid-token").
			Return(nil, domain.ErrTokenInvalid).
			Once()

		result, err :=
			svc.ValidateToken(
				context.Background(),
				"invalid-token",
			)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.False(t, result.Valid)
	})

	t.Run("user lookup fails", func(t *testing.T) {
		svc, _, _, tokenManager, _, users :=
			newTestService()

		expiresAt := time.Now().Add(time.Hour)

		claims := &domain.AccessTokenClaims{
			UserID:    42,
			ExpiresAt: expiresAt,
		}

		tokenManager.
			On("ValidateAccessToken", "valid-token").
			Return(claims, nil).
			Once()

		users.
			On("GetByID", mock.Anything, 42).
			Return(nil, domain.ErrUserNotFound).
			Once()

		result, err :=
			svc.ValidateToken(
				context.Background(),
				"valid-token",
			)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.False(t, result.Valid)
	})

	t.Run("success", func(t *testing.T) {
		svc, _, _, tokenManager, _, users :=
			newTestService()

		expiresAt := time.Now().Add(time.Hour)

		claims := &domain.AccessTokenClaims{
			// #nosec G115 -- testUser.ID is 42 which fits in uint32
			UserID:    uint32(testUser.ID),
			ExpiresAt: expiresAt,
		}

		tokenManager.
			On("ValidateAccessToken", "valid-token").
			Return(claims, nil).
			Once()

		users.
			On("GetByID", mock.Anything, testUser.ID).
			Return(testUser, nil).
			Once()

		result, err :=
			svc.ValidateToken(
				context.Background(),
				"valid-token",
			)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.True(t, result.Valid)
		require.Equal(t, testUser, result.User)
		require.Equal(t, expiresAt, result.ExpiresAt)
	})
}
