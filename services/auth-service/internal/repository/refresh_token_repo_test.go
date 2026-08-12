//go:build integration
// +build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/SteeperMold/Orbitalik/auth-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/models"
	"github.com/SteeperMold/Orbitalik/auth-service/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clearRefreshTokens(t *testing.T) {
	t.Helper()

	_, err := testDB.Exec(
		"TRUNCATE refresh_tokens RESTART IDENTITY CASCADE",
	)
	require.NoError(t, err)
}

func TestRefreshTokenRepository_Save(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	expiresAt := time.Now().Add(24 * time.Hour).UTC()

	token := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "hash-123",
		ExpiresAt: expiresAt,
	}

	err := repo.Save(ctx, token)
	require.NoError(t, err)

	var (
		id          int
		userID      int
		hash        string
		expiresAtDB time.Time
		revokedAt   *time.Time
		createdAt   time.Time
	)

	const query = `
		SELECT
			id,
			user_id,
			hash,
			expires_at,
			created_at,
			revoked_at
		FROM refresh_tokens
		WHERE hash = $1
	`

	err = testDB.
		QueryRow(query, token.Hash).
		Scan(
			&id,
			&userID,
			&hash,
			&expiresAtDB,
			&createdAt,
			&revokedAt,
		)

	require.NoError(t, err)

	assert.NotZero(t, id)
	assert.Equal(t, token.UserID, userID)
	assert.Equal(t, token.Hash, hash)
	assert.WithinDuration(t, expiresAt, expiresAtDB, time.Second)
	assert.False(t, createdAt.IsZero())
	assert.Nil(t, revokedAt)
}

func TestRefreshTokenRepository_Save_DuplicateHash(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	expiresAt := time.Now().Add(24 * time.Hour)

	first := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "duplicate-hash",
		ExpiresAt: expiresAt,
	}

	second := &models.StoredRefreshToken{
		UserID:    2,
		Hash:      "duplicate-hash",
		ExpiresAt: expiresAt,
	}

	require.NoError(t, repo.Save(ctx, first))

	err := repo.Save(ctx, second)
	require.Error(t, err)

	var count int

	err = testDB.
		QueryRow(
			"SELECT COUNT(*) FROM refresh_tokens WHERE hash = $1",
			first.Hash,
		).
		Scan(&count)

	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestRefreshTokenRepository_Save_AnyUserID(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	token := &models.StoredRefreshToken{
		UserID:    999999,
		Hash:      "hash-any-user",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	err := repo.Save(ctx, token)
	require.NoError(t, err)

	result, err := repo.GetByHash(ctx, token.Hash)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, token.UserID, result.UserID)
	assert.Equal(t, token.Hash, result.Hash)
}

func TestRefreshTokenRepository_Save_DefaultTimestamps(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	err := repo.Save(ctx, &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "timestamp-hash",
		ExpiresAt: time.Now().Add(time.Hour),
	})

	require.NoError(t, err)

	result, err := repo.GetByHash(ctx, "timestamp-hash")

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.False(t, result.CreatedAt.IsZero())
	assert.Nil(t, result.RevokedAt)
}

func TestRefreshTokenRepository_GetByHash(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	expiresAt := time.Now().Add(24 * time.Hour).UTC()

	token := &models.StoredRefreshToken{
		UserID:    42,
		Hash:      "get-by-hash",
		ExpiresAt: expiresAt,
	}

	require.NoError(t, repo.Save(ctx, token))

	result, err := repo.GetByHash(ctx, token.Hash)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.NotZero(t, result.ID)
	assert.Equal(t, token.UserID, result.UserID)
	assert.Equal(t, token.Hash, result.Hash)
	assert.WithinDuration(t, token.ExpiresAt, result.ExpiresAt, time.Second)

	assert.False(t, result.CreatedAt.IsZero())
	assert.Nil(t, result.RevokedAt)
}

func TestRefreshTokenRepository_GetByHash_NotFound(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	result, err := repo.GetByHash(ctx, "does-not-exist")

	assert.Nil(t, result)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrTokenInvalid)
}

func TestRefreshTokenRepository_GetByHash_ReturnsCorrectToken(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	tokens := []*models.StoredRefreshToken{
		{
			UserID:    1,
			Hash:      "hash-user-1",
			ExpiresAt: time.Now().Add(time.Hour),
		},
		{
			UserID:    2,
			Hash:      "hash-user-2",
			ExpiresAt: time.Now().Add(2 * time.Hour),
		},
	}

	for _, token := range tokens {
		require.NoError(t, repo.Save(ctx, token))
	}

	result, err := repo.GetByHash(ctx, "hash-user-2")

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 2, result.UserID)
	assert.Equal(t, "hash-user-2", result.Hash)
}

func TestRefreshTokenRepository_Revoke(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	token := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "revoke-hash",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	require.NoError(t, repo.Save(ctx, token))

	saved, err := repo.GetByHash(ctx, token.Hash)
	require.NoError(t, err)
	require.NotNil(t, saved)

	assert.Nil(t, saved.RevokedAt)

	err = repo.Revoke(ctx, saved.ID)
	require.NoError(t, err)

	revoked, err := repo.GetByHash(ctx, token.Hash)

	require.NoError(t, err)
	require.NotNil(t, revoked)

	assert.NotNil(t, revoked.RevokedAt)
	assert.False(t, revoked.RevokedAt.IsZero())
}

func TestRefreshTokenRepository_Revoke_NotFound(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	err := repo.Revoke(ctx, 999999)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrTokenInvalid)
}

func TestRefreshTokenRepository_Revoke_AlreadyRevoked(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	token := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "already-revoked",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	require.NoError(t, repo.Save(ctx, token))

	saved, err := repo.GetByHash(ctx, token.Hash)
	require.NoError(t, err)

	require.NoError(t, repo.Revoke(ctx, saved.ID))

	err = repo.Revoke(ctx, saved.ID)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrTokenInvalid)
}

func TestRefreshTokenRepository_Revoke_DoesNotChangeOtherTokens(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	first := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "revoke-first",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	second := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "revoke-second",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	require.NoError(t, repo.Save(ctx, first))
	require.NoError(t, repo.Save(ctx, second))

	savedFirst, err := repo.GetByHash(ctx, first.Hash)
	require.NoError(t, err)

	require.NoError(t, repo.Revoke(ctx, savedFirst.ID))

	savedSecond, err := repo.GetByHash(ctx, second.Hash)

	require.NoError(t, err)
	require.NotNil(t, savedSecond)

	assert.Nil(t, savedSecond.RevokedAt)
}

func TestRefreshTokenRepository_Rotate(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	oldToken := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "old-refresh-hash",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	require.NoError(t, repo.Save(ctx, oldToken))

	oldSaved, err := repo.GetByHash(ctx, oldToken.Hash)
	require.NoError(t, err)
	require.NotNil(t, oldSaved)

	newToken := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "new-refresh-hash",
		ExpiresAt: time.Now().Add(2 * time.Hour),
	}

	err = repo.Rotate(ctx, oldSaved.ID, newToken)
	require.NoError(t, err)

	// Old token should now be revoked.
	oldResult, err := repo.GetByHash(ctx, oldToken.Hash)

	require.NoError(t, err)
	require.NotNil(t, oldResult)

	assert.NotNil(t, oldResult.RevokedAt)
	assert.False(t, oldResult.RevokedAt.IsZero())

	// New token should exist and be active.
	newResult, err := repo.GetByHash(ctx, newToken.Hash)

	require.NoError(t, err)
	require.NotNil(t, newResult)

	assert.NotZero(t, newResult.ID)
	assert.NotEqual(t, oldResult.ID, newResult.ID)
	assert.Equal(t, newToken.UserID, newResult.UserID)
	assert.Equal(t, newToken.Hash, newResult.Hash)
	assert.Nil(t, newResult.RevokedAt)
	assert.WithinDuration(
		t,
		newToken.ExpiresAt,
		newResult.ExpiresAt,
		time.Second,
	)
}

func TestRefreshTokenRepository_Rotate_NotFound(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	newToken := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "new-token-after-failed-rotation",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	err := repo.Rotate(ctx, 999999, newToken)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrTokenInvalid)

	// The transaction should have rolled back, so the new token
	// must NOT have been inserted.
	var count int

	err = testDB.
		QueryRow(
			"SELECT COUNT(*) FROM refresh_tokens WHERE hash = $1",
			newToken.Hash,
		).
		Scan(&count)

	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestRefreshTokenRepository_Rotate_DuplicateNewHashRollsBack(
	t *testing.T,
) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	oldToken := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "rotation-old-hash",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	existingToken := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "existing-hash",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	require.NoError(t, repo.Save(ctx, oldToken))
	require.NoError(t, repo.Save(ctx, existingToken))

	oldSaved, err := repo.GetByHash(ctx, oldToken.Hash)
	require.NoError(t, err)

	newToken := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      existingToken.Hash,
		ExpiresAt: time.Now().Add(2 * time.Hour),
	}

	err = repo.Rotate(ctx, oldSaved.ID, newToken)

	require.Error(t, err)

	// The old token must NOT have been revoked because the
	// transaction should have rolled back.
	oldResult, err := repo.GetByHash(ctx, oldToken.Hash)

	require.NoError(t, err)
	require.NotNil(t, oldResult)

	assert.Nil(t, oldResult.RevokedAt)

	// There should still be exactly one token with the duplicate hash.
	var count int

	err = testDB.
		QueryRow(
			"SELECT COUNT(*) FROM refresh_tokens WHERE hash = $1",
			existingToken.Hash,
		).
		Scan(&count)

	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestRefreshTokenRepository_Rotate_DifferentUserID(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	oldToken := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "old-user-token",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	require.NoError(t, repo.Save(ctx, oldToken))

	oldSaved, err := repo.GetByHash(ctx, oldToken.Hash)
	require.NoError(t, err)

	newToken := &models.StoredRefreshToken{
		UserID:    2,
		Hash:      "new-user-token",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	err = repo.Rotate(ctx, oldSaved.ID, newToken)
	require.NoError(t, err)

	result, err := repo.GetByHash(ctx, newToken.Hash)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 2, result.UserID)
	assert.Equal(t, newToken.Hash, result.Hash)
}

func TestRefreshTokenRepository_Rotate_AlreadyRevokedOldToken(
	t *testing.T,
) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	oldToken := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "already-revoked-old",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	require.NoError(t, repo.Save(ctx, oldToken))

	oldSaved, err := repo.GetByHash(ctx, oldToken.Hash)
	require.NoError(t, err)

	require.NoError(t, repo.Revoke(ctx, oldSaved.ID))

	newToken := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "should-not-be-created",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	err = repo.Rotate(ctx, oldSaved.ID, newToken)

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrTokenInvalid)

	// Rotation must be atomic. The new token must not remain
	// after the failed update.
	var count int

	err = testDB.
		QueryRow(
			"SELECT COUNT(*) FROM refresh_tokens WHERE hash = $1",
			newToken.Hash,
		).
		Scan(&count)

	require.NoError(t, err)
	assert.Equal(t, 0, count)

	// The old token must still be revoked.
	oldResult, err := repo.GetByHash(ctx, oldToken.Hash)

	require.NoError(t, err)
	require.NotNil(t, oldResult)

	assert.NotNil(t, oldResult.RevokedAt)
}

func TestRefreshTokenRepository_FullLifecycle(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	ctx := context.Background()
	repo := repository.NewRefreshTokenRepository(testConn)

	// Create old token.
	oldToken := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "lifecycle-old",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	require.NoError(t, repo.Save(ctx, oldToken))

	// Read old token.
	oldSaved, err := repo.GetByHash(ctx, oldToken.Hash)

	require.NoError(t, err)
	require.NotNil(t, oldSaved)

	assert.Equal(t, oldToken.UserID, oldSaved.UserID)
	assert.Equal(t, oldToken.Hash, oldSaved.Hash)
	assert.Nil(t, oldSaved.RevokedAt)

	// Rotate.
	newToken := &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "lifecycle-new",
		ExpiresAt: time.Now().Add(2 * time.Hour),
	}

	err = repo.Rotate(ctx, oldSaved.ID, newToken)
	require.NoError(t, err)

	// Old token is revoked.
	oldResult, err := repo.GetByHash(ctx, oldToken.Hash)

	require.NoError(t, err)
	require.NotNil(t, oldResult)

	assert.NotNil(t, oldResult.RevokedAt)

	// New token is active.
	newResult, err := repo.GetByHash(ctx, newToken.Hash)

	require.NoError(t, err)
	require.NotNil(t, newResult)

	assert.Equal(t, newToken.UserID, newResult.UserID)
	assert.Equal(t, newToken.Hash, newResult.Hash)
	assert.Nil(t, newResult.RevokedAt)

	// Explicitly revoke new token.
	err = repo.Revoke(ctx, newResult.ID)
	require.NoError(t, err)

	finalResult, err := repo.GetByHash(ctx, newToken.Hash)

	require.NoError(t, err)
	require.NotNil(t, finalResult)

	assert.NotNil(t, finalResult.RevokedAt)
}

func TestRefreshTokenRepository_Save_ContextCancellation(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	repo := repository.NewRefreshTokenRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.Save(ctx, &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "cancelled-save",
		ExpiresAt: time.Now().Add(time.Hour),
	})

	require.Error(t, err)
}

func TestRefreshTokenRepository_GetByHash_ContextCancellation(
	t *testing.T,
) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	repo := repository.NewRefreshTokenRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := repo.GetByHash(ctx, "cancelled-get")

	assert.Nil(t, result)
	require.Error(t, err)
}

func TestRefreshTokenRepository_Revoke_ContextCancellation(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	repo := repository.NewRefreshTokenRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.Revoke(ctx, 1)

	require.Error(t, err)
}

func TestRefreshTokenRepository_Rotate_ContextCancellation(t *testing.T) {
	t.Cleanup(func() {
		clearRefreshTokens(t)
	})

	repo := repository.NewRefreshTokenRepository(testConn)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.Rotate(ctx, 1, &models.StoredRefreshToken{
		UserID:    1,
		Hash:      "cancelled-rotation",
		ExpiresAt: time.Now().Add(time.Hour),
	})

	require.Error(t, err)
}
