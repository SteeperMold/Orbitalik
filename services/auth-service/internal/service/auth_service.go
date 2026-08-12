package service

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
)

type AuthService struct {
	credentials  domain.CredentialsRepository
	tokens       domain.RefreshTokenRepository
	tokenManager domain.TokenManager
	hasher       domain.PasswordHasher
	users        domain.UserClient
}

func NewAuthService(
	credentials domain.CredentialsRepository,
	tokens domain.RefreshTokenRepository,
	tokenManager domain.TokenManager,
	hasher domain.PasswordHasher,
	users domain.UserClient,
) *AuthService {
	return &AuthService{
		credentials:  credentials,
		tokens:       tokens,
		tokenManager: tokenManager,
		hasher:       hasher,
		users:        users,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	email string,
	username string,
	password string,
) (*models.User, models.AccessToken, *models.RefreshToken, error) {

	if err := validateRegister(email, username, password); err != nil {
		return nil, "", nil, err
	}

	pwdHash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, "", nil, err
	}

	user, err := s.users.CreateUser(ctx, email, username)
	if err != nil {
		return nil, "", nil, err
	}

	err = s.credentials.Create(ctx, &models.Credentials{
		UserID:       user.ID,
		PasswordHash: pwdHash,
	})
	if err != nil {
		return nil, "", nil, err
	}

	access, err := s.tokenManager.GenerateAccessToken(user)
	if err != nil {
		return nil, "", nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, "", nil, err
	}

	err = s.tokens.Save(ctx, &models.StoredRefreshToken{
		UserID:    user.ID,
		Hash:      refreshToken.Hash,
		ExpiresAt: refreshToken.ExpiresAt,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, "", nil, err
	}

	return user, access, refreshToken, nil
}

func validateRegister(
	email string,
	username string,
	password string,
) error {
	email = strings.TrimSpace(email)
	username = strings.TrimSpace(username)

	if email == "" {
		return domain.ErrInvalidEmail
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return domain.ErrInvalidEmail
	}

	if username == "" {
		return domain.ErrInvalidUsername
	}

	usernameLength := utf8.RuneCountInString(username)

	if usernameLength < 3 || usernameLength > 32 {
		return domain.ErrInvalidUsername
	}

	if password == "" {
		return domain.ErrInvalidPassword
	}

	if len(password) < 8 {
		return domain.ErrWeakPassword
	}

	if len(password) > 128 {
		return domain.ErrPasswordTooLong
	}

	return nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (*models.User, models.AccessToken, *models.RefreshToken, error) {

	if err := validateLogin(email, password); err != nil {
		return nil, "", nil, err
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", nil, err
	}

	credentials, err := s.credentials.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, "", nil, domain.ErrInvalidCredentials
	}

	if err := s.hasher.Compare(
		credentials.PasswordHash,
		password,
	); err != nil {
		return nil, "", nil, domain.ErrInvalidCredentials
	}

	access, err := s.tokenManager.GenerateAccessToken(user)
	if err != nil {
		return nil, "", nil, err
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, "", nil, err
	}

	err = s.tokens.Save(ctx, &models.StoredRefreshToken{
		UserID:    user.ID,
		Hash:      refreshToken.Hash,
		ExpiresAt: refreshToken.ExpiresAt,
		CreatedAt: time.Now(),
	})

	if err != nil {
		return nil, "", nil, err
	}

	return user, access, refreshToken, nil
}

func validateLogin(email, password string) error {
	if email == "" {
		return domain.ErrInvalidEmail
	}

	if password == "" {
		return domain.ErrInvalidPassword
	}

	return nil
}

func (s *AuthService) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*models.User, models.AccessToken, *models.RefreshToken, error) {

	if refreshToken == "" {
		return nil, "", nil, domain.ErrTokenInvalid
	}

	hash := s.tokenManager.HashToken(refreshToken)

	oldToken, err := s.tokens.GetByHash(ctx, hash)
	if err != nil {
		return nil, "", nil, domain.ErrTokenInvalid
	}

	if time.Now().After(oldToken.ExpiresAt) {
		return nil, "", nil, domain.ErrTokenExpired
	}

	if oldToken.RevokedAt != nil {
		return nil, "", nil, domain.ErrTokenInvalid
	}

	user, err := s.users.GetByID(ctx, oldToken.UserID)
	if err != nil {
		return nil, "", nil, err
	}

	access, err := s.tokenManager.GenerateAccessToken(user)
	if err != nil {
		return nil, "", nil, err
	}

	issued, err := s.tokenManager.GenerateRefreshToken()
	if err != nil {
		return nil, "", nil, err
	}

	err = s.tokens.Rotate(ctx, oldToken.ID, &models.StoredRefreshToken{
		UserID:    user.ID,
		Hash:      issued.Hash,
		ExpiresAt: issued.ExpiresAt,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, "", nil, err
	}

	return user, access, issued, nil
}

func (s *AuthService) Logout(
	ctx context.Context,
	rawToken string,
) error {

	if rawToken == "" {
		return domain.ErrTokenInvalid
	}

	hash := s.tokenManager.HashToken(rawToken)

	token, err := s.tokens.GetByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, domain.ErrTokenInvalid) {
			// not revealing whether it existed
			return nil
		}
		// returning only infrastructure errors
		return err
	}

	if token.RevokedAt != nil {
		return nil
	}

	err = s.tokens.Revoke(ctx, token.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ValidateToken(
	ctx context.Context,
	rawToken string,
) (*domain.TokenValidationResult, error) {

	if rawToken == "" {
		return &domain.TokenValidationResult{
			Valid: false,
		}, nil
	}

	claims, err := s.tokenManager.ValidateAccessToken(rawToken)
	if err != nil {
		if errors.Is(err, domain.ErrTokenInvalid) {
			return &domain.TokenValidationResult{Valid: false}, nil
		}
		return nil, err
	}

	user, err := s.users.GetByID(ctx, int(claims.UserID))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return &domain.TokenValidationResult{Valid: false}, nil
		}
		return nil, err
	}

	return &domain.TokenValidationResult{
		Valid:     true,
		User:      user,
		ExpiresAt: claims.ExpiresAt,
	}, nil
}
