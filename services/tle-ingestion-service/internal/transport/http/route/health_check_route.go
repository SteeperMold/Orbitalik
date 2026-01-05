package route

import (
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/transport/http/handler"
	"net/http"
	"time"

	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// NewHealthCheckRoute registers the /health endpoint on the provided router.
func NewHealthCheckRoute(mux *mux.Router, db domain.DBConn, logger *zap.Logger, timeout time.Duration) {
	hs := service.NewHealthCheckService(db)
	hh := handler.NewHealthHandler(hs, logger, timeout)

	mux.HandleFunc("/health", hh.HealthCheck).Methods(http.MethodGet)
}
