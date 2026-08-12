package client

import (
	"context"
	"errors"
	"testing"

	"github.com/SteeperMold/Orbitalik/auth-service/gen/userpb"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// MockUserServiceClient is a Testify mock for userpb.UserServiceClient.
type MockUserServiceClient struct {
	mock.Mock
}

func (m *MockUserServiceClient) UpdateUser(
	ctx context.Context,
	in *userpb.UpdateUserRequest,
	opts ...grpc.CallOption,
) (*userpb.UpdateUserResponse, error) {
	panic("unimplemented")
}

func (m *MockUserServiceClient) DeleteUser(
	ctx context.Context,
	in *userpb.DeleteUserRequest,
	opts ...grpc.CallOption,
) (*userpb.DeleteUserResponse, error) {
	panic("unimplemented")
}

func (m *MockUserServiceClient) CreateUser(
	ctx context.Context,
	in *userpb.CreateUserRequest,
	opts ...grpc.CallOption,
) (*userpb.CreateUserResponse, error) {
	args := m.Called(ctx, in)

	var resp *userpb.CreateUserResponse
	if args.Get(0) != nil {
		resp = args.Get(0).(*userpb.CreateUserResponse)
	}

	return resp, args.Error(1)
}

func (m *MockUserServiceClient) GetUserByEmail(
	ctx context.Context,
	in *userpb.GetUserByEmailRequest,
	opts ...grpc.CallOption,
) (*userpb.GetUserResponse, error) {
	args := m.Called(ctx, in)

	var resp *userpb.GetUserResponse
	if args.Get(0) != nil {
		resp = args.Get(0).(*userpb.GetUserResponse)
	}

	return resp, args.Error(1)
}

//revive:disable-next-line:var-naming -- method name is required by the grpc interface
func (m *MockUserServiceClient) GetUserById(
	ctx context.Context,
	in *userpb.GetUserByIdRequest,
	opts ...grpc.CallOption,
) (*userpb.GetUserResponse, error) {
	args := m.Called(ctx, in)

	var resp *userpb.GetUserResponse
	if args.Get(0) != nil {
		resp = args.Get(0).(*userpb.GetUserResponse)
	}

	return resp, args.Error(1)
}

func newTestClient(mockClient *MockUserServiceClient) *Client {
	return &Client{
		client: mockClient,
	}
}

func TestClient_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx := context.Background()

		expectedResponse := &userpb.CreateUserResponse{
			User: &userpb.User{
				Id:       123,
				Email:    "john@example.com",
				Username: "john",
			},
		}

		mockClient.
			On(
				"CreateUser",
				ctx,
				&userpb.CreateUserRequest{
					Email:    "john@example.com",
					Username: "john",
				},
			).
			Return(expectedResponse, nil).
			Once()

		client := newTestClient(mockClient)

		user, err := client.CreateUser(
			ctx,
			"john@example.com",
			"john",
		)

		require.NoError(t, err)
		require.NotNil(t, user)

		assert.Equal(t, &models.User{
			ID:       123,
			Email:    "john@example.com",
			Username: "john",
		}, user)

		mockClient.AssertExpectations(t)
	})

	t.Run("already exists", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx := context.Background()

		mockClient.
			On(
				"CreateUser",
				ctx,
				&userpb.CreateUserRequest{
					Email:    "john@example.com",
					Username: "john",
				},
			).
			Return(
				(*userpb.CreateUserResponse)(nil),
				status.Error(codes.AlreadyExists, "user already exists"),
			).
			Once()

		client := newTestClient(mockClient)

		user, err := client.CreateUser(
			ctx,
			"john@example.com",
			"john",
		)

		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)

		mockClient.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx := context.Background()

		mockClient.
			On(
				"CreateUser",
				ctx,
				&userpb.CreateUserRequest{
					Email:    "john@example.com",
					Username: "john",
				},
			).
			Return(
				(*userpb.CreateUserResponse)(nil),
				status.Error(codes.NotFound, "user not found"),
			).
			Once()

		client := newTestClient(mockClient)

		user, err := client.CreateUser(
			ctx,
			"john@example.com",
			"john",
		)

		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)

		mockClient.AssertExpectations(t)
	})

	t.Run("unknown grpc error", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx := context.Background()

		expectedErr := status.Error(
			codes.Internal,
			"database unavailable",
		)

		mockClient.
			On(
				"CreateUser",
				ctx,
				&userpb.CreateUserRequest{
					Email:    "john@example.com",
					Username: "john",
				},
			).
			Return(
				(*userpb.CreateUserResponse)(nil),
				expectedErr,
			).
			Once()

		client := newTestClient(mockClient)

		user, err := client.CreateUser(
			ctx,
			"john@example.com",
			"john",
		)

		assert.Nil(t, user)
		assert.ErrorIs(t, err, expectedErr)

		mockClient.AssertExpectations(t)
	})

	t.Run("context is propagated", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockClient.
			On(
				"CreateUser",
				ctx,
				&userpb.CreateUserRequest{
					Email:    "john@example.com",
					Username: "john",
				},
			).
			Return(
				&userpb.CreateUserResponse{
					User: &userpb.User{
						Id:       1,
						Email:    "john@example.com",
						Username: "john",
					},
				},
				nil,
			).
			Once()

		client := newTestClient(mockClient)

		user, err := client.CreateUser(
			ctx,
			"john@example.com",
			"john",
		)

		require.NoError(t, err)
		assert.Equal(t, 1, user.ID)

		mockClient.AssertExpectations(t)
	})

	t.Run("empty values are passed through", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx := context.Background()

		mockClient.
			On(
				"CreateUser",
				ctx,
				&userpb.CreateUserRequest{
					Email:    "",
					Username: "",
				},
			).
			Return(
				&userpb.CreateUserResponse{
					User: &userpb.User{
						Id:       0,
						Email:    "",
						Username: "",
					},
				},
				nil,
			).
			Once()

		client := newTestClient(mockClient)

		user, err := client.CreateUser(ctx, "", "")

		require.NoError(t, err)
		require.NotNil(t, user)

		assert.Equal(t, &models.User{
			ID:       0,
			Email:    "",
			Username: "",
		}, user)

		mockClient.AssertExpectations(t)
	})
}

func TestClient_GetByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx := context.Background()

		expectedResponse := &userpb.GetUserResponse{
			User: &userpb.User{
				Id:       42,
				Email:    "alice@example.com",
				Username: "alice",
			},
		}

		mockClient.
			On(
				"GetUserByEmail",
				ctx,
				&userpb.GetUserByEmailRequest{
					Email: "alice@example.com",
				},
			).
			Return(expectedResponse, nil).
			Once()

		client := newTestClient(mockClient)

		user, err := client.GetByEmail(
			ctx,
			"alice@example.com",
		)

		require.NoError(t, err)
		require.NotNil(t, user)

		assert.Equal(t, &models.User{
			ID:       42,
			Email:    "alice@example.com",
			Username: "alice",
		}, user)

		mockClient.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx := context.Background()

		mockClient.
			On(
				"GetUserByEmail",
				ctx,
				&userpb.GetUserByEmailRequest{
					Email: "missing@example.com",
				},
			).
			Return(
				(*userpb.GetUserResponse)(nil),
				status.Error(codes.NotFound, "user not found"),
			).
			Once()

		client := newTestClient(mockClient)

		user, err := client.GetByEmail(
			ctx,
			"missing@example.com",
		)

		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)

		mockClient.AssertExpectations(t)
	})

	t.Run("already exists is mapped correctly", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx := context.Background()

		mockClient.
			On(
				"GetUserByEmail",
				ctx,
				&userpb.GetUserByEmailRequest{
					Email: "alice@example.com",
				},
			).
			Return(
				(*userpb.GetUserResponse)(nil),
				status.Error(codes.AlreadyExists, "already exists"),
			).
			Once()

		client := newTestClient(mockClient)

		user, err := client.GetByEmail(
			ctx,
			"alice@example.com",
		)

		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)

		mockClient.AssertExpectations(t)
	})

	t.Run("unknown error is returned unchanged", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx := context.Background()
		expectedErr := errors.New("something went wrong")

		mockClient.
			On(
				"GetUserByEmail",
				ctx,
				&userpb.GetUserByEmailRequest{
					Email: "alice@example.com",
				},
			).
			Return(
				(*userpb.GetUserResponse)(nil),
				expectedErr,
			).
			Once()

		client := newTestClient(mockClient)

		user, err := client.GetByEmail(
			ctx,
			"alice@example.com",
		)

		assert.Nil(t, user)
		assert.ErrorIs(t, err, expectedErr)

		mockClient.AssertExpectations(t)
	})
}

func TestClient_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx := context.Background()

		expectedResponse := &userpb.GetUserResponse{
			User: &userpb.User{
				Id:       100,
				Email:    "bob@example.com",
				Username: "bob",
			},
		}

		mockClient.
			On(
				"GetUserById",
				ctx,
				&userpb.GetUserByIdRequest{
					Id: 100,
				},
			).
			Return(expectedResponse, nil).
			Once()

		client := newTestClient(mockClient)

		user, err := client.GetByID(ctx, 100)

		require.NoError(t, err)
		require.NotNil(t, user)

		assert.Equal(t, &models.User{
			ID:       100,
			Email:    "bob@example.com",
			Username: "bob",
		}, user)

		mockClient.AssertExpectations(t)
	})

	t.Run("zero id", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx := context.Background()

		mockClient.
			On(
				"GetUserById",
				ctx,
				&userpb.GetUserByIdRequest{
					Id: 0,
				},
			).
			Return(
				&userpb.GetUserResponse{
					User: &userpb.User{
						Id:       0,
						Email:    "zero@example.com",
						Username: "zero",
					},
				},
				nil,
			).
			Once()

		client := newTestClient(mockClient)

		user, err := client.GetByID(ctx, 0)

		require.NoError(t, err)
		require.NotNil(t, user)

		assert.Equal(t, 0, user.ID)
		assert.Equal(t, "zero@example.com", user.Email)
		assert.Equal(t, "zero", user.Username)

		mockClient.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx := context.Background()

		mockClient.
			On(
				"GetUserById",
				ctx,
				&userpb.GetUserByIdRequest{
					Id: 999,
				},
			).
			Return(
				(*userpb.GetUserResponse)(nil),
				status.Error(codes.NotFound, "user not found"),
			).
			Once()

		client := newTestClient(mockClient)

		user, err := client.GetByID(ctx, 999)

		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)

		mockClient.AssertExpectations(t)
	})

	t.Run("already exists", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx := context.Background()

		mockClient.
			On(
				"GetUserById",
				ctx,
				&userpb.GetUserByIdRequest{
					Id: 123,
				},
			).
			Return(
				(*userpb.GetUserResponse)(nil),
				status.Error(codes.AlreadyExists, "already exists"),
			).
			Once()

		client := newTestClient(mockClient)

		user, err := client.GetByID(ctx, 123)

		assert.Nil(t, user)
		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)

		mockClient.AssertExpectations(t)
	})

	t.Run("unknown grpc error", func(t *testing.T) {
		mockClient := new(MockUserServiceClient)

		ctx := context.Background()

		expectedErr := status.Error(
			codes.Unavailable,
			"user service unavailable",
		)

		mockClient.
			On(
				"GetUserById",
				ctx,
				&userpb.GetUserByIdRequest{
					Id: 123,
				},
			).
			Return(
				(*userpb.GetUserResponse)(nil),
				expectedErr,
			).
			Once()

		client := newTestClient(mockClient)

		user, err := client.GetByID(ctx, 123)

		assert.Nil(t, user)
		assert.ErrorIs(t, err, expectedErr)

		mockClient.AssertExpectations(t)
	})
}

func TestMapError(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		err := status.Error(codes.NotFound, "user not found")

		result := mapError(err)

		assert.ErrorIs(t, result, domain.ErrUserNotFound)
	})

	t.Run("already exists", func(t *testing.T) {
		err := status.Error(codes.AlreadyExists, "user already exists")

		result := mapError(err)

		assert.ErrorIs(t, result, domain.ErrUserAlreadyExists)
	})

	t.Run("internal", func(t *testing.T) {
		err := status.Error(codes.Internal, "internal error")

		result := mapError(err)

		assert.ErrorIs(t, result, err)
	})

	t.Run("invalid argument", func(t *testing.T) {
		err := status.Error(codes.InvalidArgument, "invalid argument")

		result := mapError(err)

		assert.ErrorIs(t, result, err)
	})

	t.Run("permission denied", func(t *testing.T) {
		err := status.Error(codes.PermissionDenied, "permission denied")

		result := mapError(err)

		assert.ErrorIs(t, result, err)
	})

	t.Run("unauthenticated", func(t *testing.T) {
		err := status.Error(codes.Unauthenticated, "unauthenticated")

		result := mapError(err)

		assert.ErrorIs(t, result, err)
	})

	t.Run("plain error", func(t *testing.T) {
		err := errors.New("plain error")

		result := mapError(err)

		assert.ErrorIs(t, result, err)
	})
}

func TestClient_CloseConnection(t *testing.T) {
	conn, err := grpc.NewClient(
		"passthrough:///test",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	require.NoError(t, err)

	client := &Client{
		conn: conn,
	}

	err = client.CloseConnection()

	assert.NoError(t, err)
}
