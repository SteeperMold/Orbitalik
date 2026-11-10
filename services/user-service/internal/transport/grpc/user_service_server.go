package grpc

import (
	"context"
	"errors"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/user-service/gen/userpb"
	"github.com/SteeperMold/Orbitalik/user-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/user-service/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServiceServer struct {
	userpb.UnimplementedUserServiceServer

	service domain.UserService
	logger  applog.Logger
}

func NewUserServiceServer(s domain.UserService, logger applog.Logger) *UserServiceServer {
	return &UserServiceServer{
		service: s,
		logger:  logger,
	}
}

func (s *UserServiceServer) CreateUser(
	ctx context.Context,
	req *userpb.CreateUserRequest,
) (*userpb.CreateUserResponse, error) {

	user, err := s.service.CreateUser(ctx, &domain.CreateUserParams{
		Email:    req.GetEmail(),
		Username: req.GetUsername(),
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &userpb.CreateUserResponse{
		User: toProtoUser(user),
	}, nil
}

func (s *UserServiceServer) GetUserById(
	ctx context.Context,
	req *userpb.GetUserByIdRequest,
) (*userpb.GetUserResponse, error) {

	user, err := s.service.GetUserByID(ctx, int(req.GetId()))
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &userpb.GetUserResponse{
		User: toProtoUser(user),
	}, nil
}

func (s *UserServiceServer) GetUserByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.GetUserResponse, error) {
	user, err := s.service.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &userpb.GetUserResponse{
		User: toProtoUser(user),
	}, nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	params := domain.UpdateUserParams{
		ID: int(req.GetId()),
	}

	if req.Email != nil {
		params.Email = req.GetEmail()
	}

	if req.Username != nil {
		params.Username = req.GetUsername()
	}

	user, err := s.service.UpdateUser(ctx, &params)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &userpb.UpdateUserResponse{
		User: toProtoUser(user),
	}, nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	if err := s.service.DeleteUser(ctx, int(req.GetId())); err != nil {
		return nil, toGRPCError(err)
	}

	return &userpb.DeleteUserResponse{}, nil
}

func toProtoUser(u *models.User) *userpb.User {
	if u == nil {
		return nil
	}

	return &userpb.User{
		Id:        uint32(u.ID),
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, "User not found")

	case errors.Is(err, domain.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, "User already exists")

	case errors.Is(err, domain.ErrEmailRequired):
		return status.Error(codes.InvalidArgument, "No email")

	case errors.Is(err, domain.ErrUsernameRequired):
		return status.Error(codes.InvalidArgument, "No username")

	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
