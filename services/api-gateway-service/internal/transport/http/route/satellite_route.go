package route

import (
	"net/http"
	"time"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/service"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/transport/adapter"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/transport/http/handler"
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/gorilla/mux"
)

func NewSatellitesRoute(
	mux *mux.Router,
	connections *adapter.Connections,
	logger applog.Logger,
	timeout time.Duration,
) {
	metadataAdapter := adapter.NewMetadataAdapter(connections.Metadata)
	tleAdapter := adapter.NewTLEAdapter(connections.TLE)
	trajectoryAdapter := adapter.NewTrajectoryAdapter(connections.Trajectory)

	svc := service.NewSatelliteService(metadataAdapter, tleAdapter, trajectoryAdapter)

	sh := handler.NewSatelliteHandler(svc, logger, timeout)

	mux.HandleFunc("/satellites", sh.List).Methods(http.MethodGet)
	mux.HandleFunc("/tles", sh.ListTles).Methods(http.MethodGet)

	mux.HandleFunc("/satellites/{id}", sh.Get).Methods(http.MethodGet)

	mux.HandleFunc("/satellites/{id}/metadata", sh.GetMetadata).Methods(http.MethodGet)
	mux.HandleFunc("/satellites/{id}/frequencies", sh.GetFrequencies).Methods(http.MethodGet)
	mux.HandleFunc("/satellites/{id}/sources", sh.GetSources).Methods(http.MethodGet)

	mux.HandleFunc("/satellites/{id}/orbit", sh.GetTle).Methods(http.MethodGet)
	mux.HandleFunc("/satellites/{id}/position", sh.GetPosition).Methods(http.MethodGet)
	mux.HandleFunc("/satellites/{id}/look-angles", sh.GetLookAngles).Methods(http.MethodGet)
	mux.HandleFunc("/satellites/{id}/trajectory", sh.GetTrajectory).Methods(http.MethodGet)
	mux.HandleFunc("/satellites/{id}/observer-trajectory", sh.GetObserverTrajectory).Methods(http.MethodGet)
	mux.HandleFunc("/satellites/{id}/passes", sh.GetPasses).Methods(http.MethodGet)
	mux.HandleFunc("/satellites/{id}/next-passes", sh.GetNextPasses).Methods(http.MethodGet)
}
