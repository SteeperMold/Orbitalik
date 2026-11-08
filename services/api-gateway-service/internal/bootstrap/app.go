package bootstrap

import (
	infra "github.com/SteeperMold/Orbitalik/api-gateway-service/internal/infrastructure"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/transport/adapter"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/transport/http/route"
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
)

func StartHTTPServer(cfg *infra.Config, connections *adapter.Connections, logger applog.Logger) {
	route.Serve(cfg, connections, logger)
}
