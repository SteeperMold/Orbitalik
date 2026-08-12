package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealthCheckService_HealthCheck(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()

		mockDB := new(MockDB)

		mockDB.
			On("Ping", ctx).
			Return(nil).
			Once()

		svc := NewHealthCheckService(mockDB)

		err := svc.HealthCheck(ctx)

		require.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("database ping fails", func(t *testing.T) {
		ctx := context.Background()
		expectedErr := errors.New("database unavailable")

		mockDB := new(MockDB)

		mockDB.
			On("Ping", ctx).
			Return(expectedErr).
			Once()

		svc := NewHealthCheckService(mockDB)

		err := svc.HealthCheck(ctx)

		require.ErrorIs(t, err, expectedErr)
		mockDB.AssertExpectations(t)
	})
}
