package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/domain"
	"github.com/gorilla/mux"
)

func parseSatelliteIdentifier(r *http.Request) (domain.SatelliteIdentifier, error) {
	raw := mux.Vars(r)["id"]

	raw = strings.TrimSpace(raw)

	if raw == "" {
		return domain.SatelliteIdentifier{}, fmt.Errorf("satellite identifier is required")
	}

	if id, err := strconv.ParseInt(raw, 10, 64); err == nil {
		if id <= 0 {
			return domain.SatelliteIdentifier{}, fmt.Errorf("NORAD ID must be positive")
		}

		return domain.SatelliteIdentifier{NORADID: &id}, nil
	}

	name := raw
	return domain.SatelliteIdentifier{Name: &name}, nil
}
