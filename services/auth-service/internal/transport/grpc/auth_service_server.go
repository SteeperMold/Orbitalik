package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/SteeperMold/Orbitalik/auth-service/gen/authpb"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	"github.com/SteeperMold/Orbitalik/common/go/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServiceServer struct {
	authpb.UnimplementedAuthServiceServer

	svc    domain.AuthService
	logger log.Logger
}

func NewAuthServiceServer(svc domain.AuthService, logger log.Logger) *AuthServiceServer {
	return &AuthServiceServer{
		svc:    svc,
		logger: logger,
	}
}

func (a *AuthServiceServer) Register(
	ctx context.Context,
	in *authpb.RegisterRequest,
) (*authpb.AuthResponse, error) {

	user, access, refresh, err := a.svc.Register(ctx, in.Email, in.Username, in.Password)
	if err != nil {
		return nil, a.mapErr(err)
	}

	return buildAuthResponse(user, access, refresh), nil
}

func (a *AuthServiceServer) Login(
	ctx context.Context,
	in *authpb.LoginRequest,
) (*authpb.AuthResponse, error) {

	user, access, refresh, err := a.svc.Login(ctx, in.Email, in.Password)
	if err != nil {
		return nil, a.mapErr(err)
	}

	return buildAuthResponse(user, access, refresh), nil
}

func (a *AuthServiceServer) RefreshToken(
	ctx context.Context,
	in *authpb.RefreshTokenRequest,
) (*authpb.AuthResponse, error) {

	user, access, refresh, err := a.svc.RefreshToken(ctx, in.RefreshToken)
	if err != nil {
		return nil, a.mapErr(err)
	}

	return buildAuthResponse(user, access, refresh), nil
}

func (a *AuthServiceServer) Logout(
	ctx context.Context,
	in *authpb.LogoutRequest,
) (*authpb.LogoutResponse, error) {

	err := a.svc.Logout(ctx, in.RefreshToken)
	if err != nil {
		return nil, a.mapErr(err)
	}

	return &authpb.LogoutResponse{}, nil
}

func (a *AuthServiceServer) ValidateToken(
	ctx context.Context,
	in *authpb.ValidateTokenRequest,
) (*authpb.ValidateTokenResponse, error) {

	validationResult, err := a.svc.ValidateToken(ctx, in.AccessToken)
	if err != nil {
		return nil, a.mapErr(err)
	}

	response := &authpb.ValidateTokenResponse{
		Valid: validationResult.Valid,
	}

	if !validationResult.Valid {
		return response, nil
	}

	response.User = userToProto(validationResult.User)
	response.ExpiresAt = timestamppb.New(validationResult.ExpiresAt)

	return response, nil
}

func userToProto(user *models.User) *authpb.AuthUser {
	return &authpb.AuthUser{
		// #nosec G115 -- user.ID is a postgres INTEGER and fits in uint32
		Id:       uint32(user.ID),
		Email:    user.Email,
		Username: user.Username,
	}
}

func tokensToProto(
	access models.AccessToken,
	refresh *models.RefreshToken,
) *authpb.TokenPair {

	return &authpb.TokenPair{
		AccessToken:      string(access),
		RefreshToken:     refresh.Value,
		ExpiresInSeconds: int64(time.Until(refresh.ExpiresAt).Seconds()),
	}
}

func buildAuthResponse(
	user *models.User,
	access models.AccessToken,
	refresh *models.RefreshToken,
) *authpb.AuthResponse {

	return &authpb.AuthResponse{
		User:   userToProto(user),
		Tokens: tokensToProto(access, refresh),
	}
}

func (a *AuthServiceServer) mapErr(err error) error {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")

	case errors.Is(err, domain.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, "user already exists")

	case errors.Is(err, domain.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, "invalid credentials")

	case errors.Is(err, domain.ErrTokenInvalid):
		return status.Error(codes.Unauthenticated, "invalid token")

	case errors.Is(err, domain.ErrTokenExpired):
		return status.Error(codes.Unauthenticated, "token expired")

	case errors.Is(err, domain.ErrInvalidEmail):
		return status.Error(codes.InvalidArgument, "invalid email")

	case errors.Is(err, domain.ErrInvalidUsername):
		return status.Error(codes.InvalidArgument, "invalid username")

	case errors.Is(err, domain.ErrInvalidPassword):
		return status.Error(codes.InvalidArgument, "invalid password")

	case errors.Is(err, domain.ErrWeakPassword):
		return status.Error(codes.InvalidArgument, "password is too weak")

	case errors.Is(err, domain.ErrPasswordTooLong):
		return status.Error(codes.InvalidArgument, "password is too long")

	default:
		a.logger.Error("internal server error", log.NewErrorField(err))
		return status.Error(codes.Internal, "internal server error")
	}
}
