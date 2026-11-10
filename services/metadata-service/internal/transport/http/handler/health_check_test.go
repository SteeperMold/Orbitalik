package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	commonlog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/transport/http/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHealthCheckService struct {
	mock.Mock
}

func (m *MockHealthCheckService) HealthCheck(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Info(msg string, fields ...commonlog.Field) {
	m.Called(msg, fields)
}

func (m *mockLogger) Sync() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockLogger) Error(msg string, fields ...commonlog.Field) {
	m.Called(msg, fields)
}

func TestHealthCheckHandler_HealthCheck(t *testing.T) {
	tests := []struct {
		name       string
		mockSetup  func(m *MockHealthCheckService)
		wantStatus int
	}{
		{
			name: "healthy",
			mockSetup: func(m *MockHealthCheckService) {
				m.
					On("HealthCheck", mock.Anything).
					Return(nil).
					Once()
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "unhealthy",
			mockSetup: func(m *MockHealthCheckService) {
				m.
					On("HealthCheck", mock.Anything).
					Return(assert.AnError).
					Once()
			},
			wantStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := new(MockHealthCheckService)
			tc.mockSetup(mockSvc)

			logger := &mockLogger{}
			h := handler.NewHealthHandler(mockSvc, logger, 2*time.Second)

			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rr := httptest.NewRecorder()

			h.HealthCheck(rr, req)

			resp := rr.Result()
			defer func() {
				_ = resp.Body.Close()
			}()

			assert.Equal(t, tc.wantStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}
