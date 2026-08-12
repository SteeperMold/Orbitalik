package client

import (
	"context"

	"github.com/SteeperMold/Orbitalik/auth-service/gen/userpb"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
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

	return &models.User{
		ID:       int(resp.User.Id),
		Email:    resp.User.Email,
		Username: resp.User.Username,
	}, nil
}

func (c *Client) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	resp, err := c.client.GetUserByEmail(
		ctx,
		&userpb.GetUserByEmailRequest{Email: email},
	)

	if err != nil {
		return nil, mapError(err)
	}

	return &models.User{
		ID:       int(resp.User.Id),
		Email:    resp.User.Email,
		Username: resp.User.Username,
	}, nil
}

func (c *Client) GetByID(ctx context.Context, id int) (*models.User, error) {
	resp, err := c.client.GetUserById(
		ctx,
		&userpb.GetUserByIdRequest{
			// #nosec G115 -- user.ID is a postgres INTEGER and fits in uint32
			Id: uint32(id),
		},
	)

	if err != nil {
		return nil, mapError(err)
	}

	return &models.User{
		ID:       int(resp.User.Id),
		Email:    resp.User.Email,
		Username: resp.User.Username,
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
