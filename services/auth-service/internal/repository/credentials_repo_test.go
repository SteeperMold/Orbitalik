//go:build integration
// +build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clearCredentials(t *testing.T) {
	t.Helper()

	_, err := testDB.Exec(
		"TRUNCATE credentials RESTART IDENTITY CASCADE",
	)
	require.NoError(t, err)
}

func TestCredentialsRepository_Create(t *testing.T) {
	t.Cleanup(func() {
		clearCredentials(t)
	})

	ctx := context.Background()
	repo := repository.NewCredentialsRepository(testConn)

	cred := &models.Credentials{
		UserID:       1,
		PasswordHash: "hashed-password",
	}

	err := repo.Create(ctx, cred)
	require.NoError(t, err)

	var (
		userID       int
		passwordHash string
	)

	const query = `
		SELECT user_id, password_hash
		FROM credentials
		WHERE user_id = $1
	`

	err = testDB.
		QueryRow(query, cred.UserID).
		Scan(&userID, &passwordHash)

	require.NoError(t, err)

	assert.Equal(t, cred.UserID, userID)
	assert.Equal(t, cred.PasswordHash, passwordHash)
}

func TestCredentialsRepository_Create_DuplicateUserID(t *testing.T) {
	t.Cleanup(func() {
		clearCredentials(t)
	})

	ctx := context.Background()
	repo := repository.NewCredentialsRepository(testConn)

	cred := &models.Credentials{
		UserID:       1,
		PasswordHash: "first-hash",
	}

	require.NoError(t, repo.Create(ctx, cred))

	err := repo.Create(ctx, &models.Credentials{
		UserID:       1,
		PasswordHash: "second-hash",
	})

	require.Error(t, err)

	var passwordHash string

	const query = `
		SELECT password_hash
		FROM credentials
		WHERE user_id = $1
	`

	err = testDB.QueryRow(query, 1).Scan(&passwordHash)

	require.NoError(t, err)
	assert.Equal(t, "first-hash", passwordHash)
}

func TestCredentialsRepository_Create_AnyUserID(t *testing.T) {
	t.Cleanup(func() {
		clearCredentials(t)
	})

	ctx := context.Background()
	repo := repository.NewCredentialsRepository(testConn)

	cred := &models.Credentials{
		UserID:       999999,
		PasswordHash: "hashed-password",
	}

	err := repo.Create(ctx, cred)

	require.NoError(t, err)

	result, err := repo.GetByUserID(ctx, cred.UserID)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, cred.UserID, result.UserID)
	assert.Equal(t, cred.PasswordHash, result.PasswordHash)
}

func TestCredentialsRepository_Create_DefaultTimestamps(t *testing.T) {
	t.Cleanup(func() {
		clearCredentials(t)
	})

	ctx := context.Background()
	repo := repository.NewCredentialsRepository(testConn)

	err := repo.Create(ctx, &models.Credentials{
		UserID:       1,
		PasswordHash: "hashed-password",
	})

	require.NoError(t, err)

	result, err := repo.GetByUserID(ctx, 1)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.False(t, result.CreatedAt.IsZero())
	assert.False(t, result.UpdatedAt.IsZero())
}

func TestCredentialsRepository_GetByUserID(t *testing.T) {
	t.Cleanup(func() {
		clearCredentials(t)
	})

	ctx := context.Background()
	repo := repository.NewCredentialsRepository(testConn)

	cred := &models.Credentials{
		UserID:       1,
		PasswordHash: "hashed-password",
	}

	require.NoError(t, repo.Create(ctx, cred))

	result, err := repo.GetByUserID(ctx, cred.UserID)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, cred.UserID, result.UserID)
	assert.Equal(t, cred.PasswordHash, result.PasswordHash)

	assert.False(t, result.CreatedAt.IsZero())
	assert.False(t, result.UpdatedAt.IsZero())
}

func TestCredentialsRepository_GetByUserID_NotFound(t *testing.T) {
	t.Cleanup(func() {
		clearCredentials(t)
	})

	ctx := context.Background()
	repo := repository.NewCredentialsRepository(testConn)

	result, err := repo.GetByUserID(ctx, 999999)

	assert.Nil(t, result)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUserNotFound)
}

func TestCredentialsRepository_GetByUserID_ReturnsCorrectUser(t *testing.T) {
	t.Cleanup(func() {
		clearCredentials(t)
	})

	ctx := context.Background()
	repo := repository.NewCredentialsRepository(testConn)

	// #nosec G101 -- test fixture values, not real credentials
	credentials := []*models.Credentials{
		{
			UserID:       1,
			PasswordHash: "hash-user-1",
		},
		{
			UserID:       2,
			PasswordHash: "hash-user-2",
		},
	}

	for _, cred := range credentials {
		require.NoError(t, repo.Create(ctx, cred))
	}

	result, err := repo.GetByUserID(ctx, 2)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 2, result.UserID)
	assert.Equal(t, "hash-user-2", result.PasswordHash)
}

func TestCredentialsRepository_UpdatePasswordHash(t *testing.T) {
	t.Cleanup(func() {
		clearCredentials(t)
	})

	ctx := context.Background()
	repo := repository.NewCredentialsRepository(testConn)

	cred := &models.Credentials{
		UserID:       1,
		PasswordHash: "old-password-hash",
	}

	require.NoError(t, repo.Create(ctx, cred))

	before, err := repo.GetByUserID(ctx, cred.UserID)
	require.NoError(t, err)

	err = repo.UpdatePasswordHash(
		ctx,
		cred.UserID,
		"new-password-hash",
	)

	require.NoError(t, err)

	after, err := repo.GetByUserID(ctx, cred.UserID)

	require.NoError(t, err)
	require.NotNil(t, after)

	assert.Equal(t, cred.UserID, after.UserID)
	assert.Equal(t, "new-password-hash", after.PasswordHash)

	// created_at should not change.
	assert.Equal(t, before.CreatedAt, after.CreatedAt)

	assert.False(t, after.UpdatedAt.IsZero())
}

func TestCredentialsRepository_UpdatePasswordHash_NotFound(t *testing.T) {
	t.Cleanup(func() {
		clearCredentials(t)
	})

	ctx := context.Background()
	repo := repository.NewCredentialsRepository(testConn)

	err := repo.UpdatePasswordHash(
		ctx,
		999999,
		"new-password-hash",
	)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUserNotFound)
}

func TestCredentialsRepository_UpdatePasswordHash_CanUpdateMultipleTimes(t *testing.T) {
	t.Cleanup(func() {
		clearCredentials(t)
	})

	ctx := context.Background()
	repo := repository.NewCredentialsRepository(testConn)

	require.NoError(t, repo.Create(ctx, &models.Credentials{
		UserID:       1,
		PasswordHash: "hash-1",
	}))

	require.NoError(
		t,
		repo.UpdatePasswordHash(ctx, 1, "hash-2"),
	)

	require.NoError(
		t,
		repo.UpdatePasswordHash(ctx, 1, "hash-3"),
	)

	result, err := repo.GetByUserID(ctx, 1)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, result.UserID)
	assert.Equal(t, "hash-3", result.PasswordHash)
}

func TestCredentialsRepository_FullLifecycle(t *testing.T) {
	t.Cleanup(func() {
		clearCredentials(t)
	})

	ctx := context.Background()
	repo := repository.NewCredentialsRepository(testConn)

	userID := 1

	// Create.
	err := repo.Create(ctx, &models.Credentials{
		UserID:       userID,
		PasswordHash: "initial-hash",
	})
	require.NoError(t, err)

	// Read.
	created, err := repo.GetByUserID(ctx, userID)

	require.NoError(t, err)
	require.NotNil(t, created)

	assert.Equal(t, userID, created.UserID)
	assert.Equal(t, "initial-hash", created.PasswordHash)

	createdAt := created.CreatedAt

	// Update.
	err = repo.UpdatePasswordHash(
		ctx,
		userID,
		"updated-hash",
	)
	require.NoError(t, err)

	updated, err := repo.GetByUserID(ctx, userID)

	require.NoError(t, err)
	require.NotNil(t, updated)

	assert.Equal(t, userID, updated.UserID)
	assert.Equal(t, "updated-hash", updated.PasswordHash)

	assert.Equal(t, createdAt, updated.CreatedAt)
}

func TestCredentialsRepository_ContextCancellation(t *testing.T) {
	t.Cleanup(func() {
		clearCredentials(t)
	})

	repo := repository.NewCredentialsRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.Create(ctx, &models.Credentials{
		UserID:       1,
		PasswordHash: "hash",
	})

	require.Error(t, err)
}

func TestCredentialsRepository_GetByUserID_ContextCancellation(t *testing.T) {
	t.Cleanup(func() {
		clearCredentials(t)
	})

	repo := repository.NewCredentialsRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := repo.GetByUserID(ctx, 1)

	assert.Nil(t, result)
	require.Error(t, err)
}

func TestCredentialsRepository_UpdatePasswordHash_ContextCancellation(t *testing.T) {
	t.Cleanup(func() {
		clearCredentials(t)
	})

	repo := repository.NewCredentialsRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.UpdatePasswordHash(ctx, 1, "hash")

	require.Error(t, err)
}
