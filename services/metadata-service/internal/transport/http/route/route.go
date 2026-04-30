package route

import (
	stdlog "log"
	"net/http"

	"github.com/SteeperMold/Orbitalik/common/go/db"
	"github.com/SteeperMold/Orbitalik/common/go/http/middleware"
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/infrastructure"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Serve configures and starts the HTTP server with routing and middleware.
func Serve(cfg *infrastructure.Config, db db.Conn, logger applog.Logger) {
	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware(logger))

	r.Handle("/metrics", promhttp.Handler())
	NewHealthCheckRoute(r, db, logger, cfg.ContextTimeout)

	stdlog.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, r))
}
