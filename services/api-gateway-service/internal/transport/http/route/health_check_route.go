package route

import (
	"net/http"
	"time"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/service"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/transport/http/handler"
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/gorilla/mux"
)

func NewHealthCheckRoute(mux *mux.Router, logger applog.Logger, timeout time.Duration) {
	hs := service.NewHealthCheckService()
	hh := handler.NewHealthHandler(hs, logger, timeout)

	mux.HandleFunc("/health", hh.HealthCheck).Methods(http.MethodGet)
}
