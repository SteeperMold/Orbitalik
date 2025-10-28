package route

import (
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/transport/http/middleware"
	"go.uber.org/zap"
	"log"
	"net/http"

	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/infrastructure"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Serve configures and starts the HTTP server with routing and middleware.
func Serve(cfg *infrastructure.Config, db domain.DBConn, logger *zap.Logger) {
	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware(logger))

	r.Handle("/metrics", promhttp.Handler())
	NewHealthCheckRoute(r, db, logger, cfg.ContextTimeout)

	log.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, r))
}
