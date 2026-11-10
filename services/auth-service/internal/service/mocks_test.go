package service

import (
	"context"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	"github.com/SteeperMold/Orbitalik/common/go/db"
	"github.com/stretchr/testify/mock"
)

type MockCredentialsRepository struct {
	mock.Mock
}

func (m *MockCredentialsRepository) Create(
	ctx context.Context,
	cred *models.Credentials,
) error {

	args := m.Called(ctx, cred)
	return args.Error(0)
}

func (m *MockCredentialsRepository) GetByUserID(
	ctx context.Context,
	userID int,
) (*models.Credentials, error) {

	args := m.Called(ctx, userID)

	var credentials *models.Credentials
	if args.Get(0) != nil {
		credentials = args.Get(0).(*models.Credentials)
	}

	return credentials, args.Error(1)
}

func (m *MockCredentialsRepository) UpdatePasswordHash(
	ctx context.Context,
	userID int,
	hash string,
) error {

	args := m.Called(ctx, userID, hash)
	return args.Error(0)
}

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Save(
	ctx context.Context,
	token *models.StoredRefreshToken,
) error {

	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetByHash(
	ctx context.Context,
	hash string,
) (*models.StoredRefreshToken, error) {

	args := m.Called(ctx, hash)

	var token *models.StoredRefreshToken
	if args.Get(0) != nil {
		token = args.Get(0).(*models.StoredRefreshToken)
	}

	return token, args.Error(1)
}

func (m *MockRefreshTokenRepository) Rotate(
	ctx context.Context,
	oldID int,
	newToken *models.StoredRefreshToken,
) error {

	args := m.Called(ctx, oldID, newToken)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) Revoke(
	ctx context.Context,
	id int,
) error {

	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockTokenManager struct {
	mock.Mock
}

func (m *MockTokenManager) GenerateAccessToken(user *models.User) (models.AccessToken, error) {
	args := m.Called(user)

	var token models.AccessToken
	if args.Get(0) != nil {
		token = args.Get(0).(models.AccessToken)
	}

	return token, args.Error(1)
}

func (m *MockTokenManager) GenerateRefreshToken() (*models.RefreshToken, error) {
	args := m.Called()

	var token *models.RefreshToken
	if args.Get(0) != nil {
		token = args.Get(0).(*models.RefreshToken)
	}

	return token, args.Error(1)
}

func (m *MockTokenManager) HashToken(token string) string {
	args := m.Called(token)
	return args.String(0)
}

func (m *MockTokenManager) ValidateAccessToken(token string) (*domain.AccessTokenClaims, error) {
	args := m.Called(token)

	var claims *domain.AccessTokenClaims
	if args.Get(0) != nil {
		claims = args.Get(0).(*domain.AccessTokenClaims)
	}

	return claims, args.Error(1)
}

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Compare(passwordHash string, password string) error {
	args := m.Called(passwordHash, password)
	return args.Error(0)
}

type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) CreateUser(
	ctx context.Context,
	email string,
	username string,
) (*models.User, error) {

	args := m.Called(ctx, email, username)

	var user *models.User
	if args.Get(0) != nil {
		user = args.Get(0).(*models.User)
	}

	return user, args.Error(1)
}

func (m *MockUserClient) GetByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {

	args := m.Called(ctx, email)

	var user *models.User
	if args.Get(0) != nil {
		user = args.Get(0).(*models.User)
	}

	return user, args.Error(1)
}

func (m *MockUserClient) GetByID(
	ctx context.Context,
	userID int,
) (*models.User, error) {

	args := m.Called(ctx, userID)

	var user *models.User
	if args.Get(0) != nil {
		user = args.Get(0).(*models.User)
	}

	return user, args.Error(1)
}

func newTestService() (
	*AuthService,
	*MockCredentialsRepository,
	*MockRefreshTokenRepository,
	*MockTokenManager,
	*MockPasswordHasher,
	*MockUserClient,
) {
	credentials := new(MockCredentialsRepository)
	tokens := new(MockRefreshTokenRepository)
	tokenManager := new(MockTokenManager)
	hasher := new(MockPasswordHasher)
	users := new(MockUserClient)

	service := NewAuthService(
		credentials,
		tokens,
		tokenManager,
		hasher,
		users,
	)

	return service, credentials, tokens, tokenManager, hasher, users
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Query(
	ctx context.Context,
	query string,
	args ...any,
) (db.Rows, error) {
	ret := m.Called(ctx, query, args)
	return ret.Get(0).(db.Rows), ret.Error(1)
}

func (m *MockDB) QueryRow(
	ctx context.Context,
	query string,
	args ...any,
) db.Row {
	ret := m.Called(ctx, query, args)
	return ret.Get(0).(db.Row)
}

func (m *MockDB) Exec(
	ctx context.Context,
	query string,
	args ...any,
) (db.Result, error) {
	ret := m.Called(ctx, query, args)
	return ret.Get(0).(db.Result), ret.Error(1)
}

func (m *MockDB) Begin(ctx context.Context) (db.Tx, error) {
	ret := m.Called(ctx)
	return ret.Get(0).(db.Tx), ret.Error(1)
}

func (m *MockDB) Ping(ctx context.Context) error {
	ret := m.Called(ctx)
	return ret.Error(0)
}

func (m *MockDB) Close() {
	m.Called()
}
