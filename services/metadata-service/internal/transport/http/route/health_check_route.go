package route

import (
	"net/http"
	"time"

	"github.com/SteeperMold/Orbitalik/common/go/db"
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/service"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/transport/http/handler"
	"github.com/gorilla/mux"
)

func NewHealthCheckRoute(mux *mux.Router, db db.Conn, logger applog.Logger, timeout time.Duration) {
	hs := service.NewHealthCheckService(db)
	hh := handler.NewHealthHandler(hs, logger, timeout)

	mux.HandleFunc("/health", hh.HealthCheck).Methods(http.MethodGet)
}
