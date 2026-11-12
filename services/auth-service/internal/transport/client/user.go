package client

import (
	"context"
	"fmt"

	"github.com/SteeperMold/Orbitalik/auth-service/gen/userpb"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn   *grpc.ClientConn
	client userpb.UserServiceClient
}

func NewUserClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: userpb.NewUserServiceClient(conn),
	}, nil
}

func (c *Client) CloseConnection() error {
	return c.conn.Close()
}

func (c *Client) CreateUser(ctx context.Context, email, username string) (*models.User, error) {
	resp, err := c.client.CreateUser(
		ctx,
		&userpb.CreateUserRequest{
			Email:    email,
			Username: username,
		},
	)

	if err != nil {
		return nil, mapError(err)
	}

	return toDomainUser(resp.GetUser())
}

func (c *Client) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	resp, err := c.client.GetUserByEmail(
		ctx,
		&userpb.GetUserByEmailRequest{Email: email},
	)

	if err != nil {
		return nil, mapError(err)
	}

	return toDomainUser(resp.GetUser())
}

func (c *Client) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	resp, err := c.client.GetUserById(
		ctx,
		&userpb.GetUserByIdRequest{
			Id: id.String(),
		},
	)

	if err != nil {
		return nil, mapError(err)
	}

	return toDomainUser(resp.GetUser())
}

func toDomainUser(user *userpb.User) (*models.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user service returned an empty user")
	}

	id, err := uuid.Parse(user.GetId())
	if err != nil {
		return nil, fmt.Errorf("user service returned an invalid user id: %w", err)
	}

	return &models.User{
		ID:       id,
		Email:    user.GetEmail(),
		Username: user.GetUsername(),
	}, nil
}

func mapError(err error) error {
	switch status.Code(err) {

	case codes.NotFound:
		return domain.ErrUserNotFound

	case codes.AlreadyExists:
		return domain.ErrUserAlreadyExists

	default:
		return err
	}
}
