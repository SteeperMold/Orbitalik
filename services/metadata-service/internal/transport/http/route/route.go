package route

import (
	stdlog "log"
	"net/http"
	"time"

	"github.com/SteeperMold/Orbitalik/common/go/db"
	"github.com/SteeperMold/Orbitalik/common/go/http/middleware"
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/infrastructure"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(cfg *infrastructure.Config, db db.Conn, logger applog.Logger) {
	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware(logger))

	r.Handle("/metrics", promhttp.Handler())
	NewHealthCheckRoute(r, db, logger, cfg.ContextTimeout)

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	stdlog.Fatal(server.ListenAndServe())
}
